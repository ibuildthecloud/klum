package github

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	klum "github.com/jadolg/klum/pkg/apis/klum.cattle.io/v1alpha1"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

func UploadKubeconfig(userSync *klum.UserSyncGithub, kubeconfig *klum.Kubeconfig, cfg Config) error {
	githubSync := userSync.Spec.Github
	if err := githubSync.Validate(); err != nil {
		return err
	}

	kubeconfigYAML, err := toYAMLString(kubeconfig.Spec)
	if err != nil {
		return err
	}

	upToDate, hash := isSecretUpToDate(userSync, kubeconfigYAML)

	if upToDate {
		return nil
	}

	ctx := context.Background()
	time.Sleep(time.Second) // Calling GitHub continuously creates problems. This adds a buffer so all operations succeed.

	log.WithFields(log.Fields{
		"secret": githubSync.SecretName,
		"user":   userSync.Spec.User,
		"repo":   fmt.Sprintf("%s/%s", githubSync.Owner, githubSync.Repository),
		"env":    githubSync.Environment,
	}).Info("Adding secret")

	client, err := newGithubClient(cfg, githubSync.Owner, githubSync.Repository)
	if err != nil {
		return err
	}

	if githubSync.Environment == "" {
		err = createRepositorySecret(
			ctx,
			client,
			&githubSync,
			kubeconfigYAML,
		)
	} else {
		err = createRepositoryEnvSecret(
			ctx,
			client,
			&githubSync,
			kubeconfigYAML,
		)
	}

	if err == nil {
		userSync.Annotations["klum.cattle.io/lastest.upload.github"] = hash
	}

	return err
}

func isSecretUpToDate(userSync *klum.UserSyncGithub, kubeconfigYAML []byte) (bool, string) {
	h := sha256.New()
	h.Write(kubeconfigYAML)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	latestKubeconfigUploaded, present := userSync.Annotations["klum.cattle.io/lastest.upload.github"]
	if present && latestKubeconfigUploaded == hash {
		return true, hash
	}

	return false, hash
}

func DeleteKubeconfig(userSync *klum.UserSyncGithub, cfg Config) error {
	githubSync := userSync.Spec.Github
	if err := githubSync.Validate(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"secret": githubSync.SecretName,
		"user":   userSync.Spec.User,
		"repo":   fmt.Sprintf("%s/%s", githubSync.Owner, githubSync.Repository),
		"env":    githubSync.Environment,
	}).Info("Deleting secret")

	client, err := newGithubClient(cfg, githubSync.Owner, githubSync.Repository)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if githubSync.Environment == "" {
		return deleteRepositorySecret(
			ctx,
			client,
			&githubSync,
		)
	} else {
		return deleteRepositoryEnvSecret(
			ctx,
			client,
			&githubSync,
		)
	}
}

func toYAMLString(x interface{}) ([]byte, error) {
	b, err := yaml.Marshal(x)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
