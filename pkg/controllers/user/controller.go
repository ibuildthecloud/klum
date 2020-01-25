package user

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	klum "github.com/ibuildthecloud/klum/pkg/apis/klum.cattle.io/v1alpha1"
	"github.com/ibuildthecloud/klum/pkg/generated/controllers/klum.cattle.io/v1alpha1"
	v1controller "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	rbaccontroller "github.com/rancher/wrangler-api/pkg/generated/controllers/rbac/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/generic"
	name2 "github.com/rancher/wrangler/pkg/name"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type Config struct {
	Namespace          string
	ContextName        string
	Server             string
	CA                 string
	DefaultClusterRole string
}

func Register(ctx context.Context,
	cfg Config,
	apply apply.Apply,
	serviceAccount v1controller.ServiceAccountController,
	crb rbaccontroller.ClusterRoleBindingController,
	secrets v1controller.SecretController,
	kconfig v1alpha1.KubeconfigController,
	user v1alpha1.UserController) {

	h := &handler{
		cfg:             cfg,
		apply:           apply.WithCacheTypes(kconfig),
		serviceAccounts: serviceAccount.Cache(),
	}

	v1alpha1.RegisterUserGeneratingHandler(ctx,
		user,
		apply.WithCacheTypes(serviceAccount,
			crb),
		"",
		"klum-user",
		h.OnUserChange,
		&generic.GeneratingHandlerOptions{
			AllowClusterScoped: true,
		})

	secrets.OnChange(ctx, "klum-secret", h.OnSecretChange)
}

type handler struct {
	cfg             Config
	apply           apply.Apply
	serviceAccounts v1controller.ServiceAccountCache
}

func (h *handler) OnUserChange(user *klum.User, status klum.UserStatus) ([]runtime.Object, klum.UserStatus, error) {
	if user.Spec.Enabled != nil && !*user.Spec.Enabled {
		status = setReady(status, false)
		return nil, status, nil
	}

	objs := []runtime.Object{
		&v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      user.Name,
				Namespace: h.cfg.Namespace,
				Annotations: map[string]string{
					"klum.cattle.io/user": user.Name,
				},
			},
		},
	}
	objs = append(objs, h.getRoles(user)...)

	return objs, setReady(status, true), nil
}

func (h *handler) getRoles(user *klum.User) []runtime.Object {
	subjects := []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      user.Name,
			Namespace: h.cfg.Namespace,
		},
	}

	if len(user.Spec.ClusterRoles) == 0 && len(user.Spec.Roles) == 0 {
		if h.cfg.DefaultClusterRole == "" {
			return nil
		}
		return []runtime.Object{
			&rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: name(user.Name, "", h.cfg.DefaultClusterRole, ""),
				},
				Subjects: subjects,
				RoleRef: rbacv1.RoleRef{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "ClusterRole",
					Name:     h.cfg.DefaultClusterRole,
				},
			},
		}
	}

	var objs []runtime.Object

	for _, clusterRole := range user.Spec.ClusterRoles {
		objs = append(objs, &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: name(user.Name, "", clusterRole, ""),
			},
			Subjects: subjects,
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     clusterRole,
			},
		})
	}

	for _, role := range user.Spec.Roles {
		if role.Namespace == "" ||
			role.Role == "" && role.ClusterRole == "" {
			continue
		}

		rb := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name(user.Name, role.Namespace, role.ClusterRole, role.Role),
				Namespace: role.Namespace,
			},
			Subjects: subjects,
		}

		if role.Role != "" {
			rb.RoleRef = rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     role.Role,
			}
		} else {
			rb.RoleRef = rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     role.ClusterRole,
			}
		}

		objs = append(objs, rb)
	}

	return objs
}

func name(user, namespace, clusterRole, role string) string {
	// this is so we don't get conflics
	suffix := md5.Sum([]byte(fmt.Sprintf("%s/%s/%s/%s", user, namespace, clusterRole, role)))
	return name2.SafeConcatName("klum", user, role, hex.EncodeToString(suffix[:])[:8])
}

func (h *handler) OnSecretChange(key string, secret *v1.Secret) (*v1.Secret, error) {
	if secret == nil {
		return nil, nil
	}

	if secret.Type != v1.SecretTypeServiceAccountToken {
		return secret, nil
	}

	sa, err := h.serviceAccounts.Get(secret.Namespace, secret.Annotations["kubernetes.io/service-account.name"])
	if errors.IsNotFound(err) {
		return secret, nil
	} else if err != nil {
		return secret, err
	}

	if sa.UID != types.UID(secret.Annotations["kubernetes.io/service-account.uid"]) {
		return secret, nil
	}

	userName := sa.Annotations["klum.cattle.io/user"]
	if userName == "" {
		return secret, nil
	}

	ca := h.cfg.CA
	if ca == "" {
		ca = base64.StdEncoding.EncodeToString(secret.Data["ca.crt"])
	}
	token := string(secret.Data["token"])

	return secret, h.apply.
		WithOwner(secret).
		WithSetOwnerReference(true, false).
		ApplyObjects(&klum.Kubeconfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: userName,
			},
			Spec: klum.KubeconfigSpec{
				Clusters: []klum.NamedCluster{
					{
						Name: h.cfg.ContextName,
						Cluster: klum.Cluster{
							Server:                   h.cfg.Server,
							CertificateAuthorityData: ca,
						},
					},
				},
				AuthInfos: []klum.NamedAuthInfo{
					{
						Name: h.cfg.ContextName,
						AuthInfo: klum.AuthInfo{
							Token: token,
						},
					},
				},
				Contexts: []klum.NamedContext{
					{
						Name: h.cfg.ContextName,
						Context: klum.Context{
							Cluster:  h.cfg.ContextName,
							AuthInfo: h.cfg.ContextName,
						},
					},
				},
				CurrentContext: h.cfg.ContextName,
			},
		})
}

func setReady(status klum.UserStatus, ready bool) klum.UserStatus {
	// dumb hack to set condition, should really make this easier
	user := &klum.User{Status: status}
	klum.UserReadyCondition.SetStatusBool(user, ready)
	return user.Status
}
