//go:generate go run pkg/codegen/cleanup/main.go
//go:generate go run pkg/codegen/main.go

package main

import (
	"fmt"
	"os"

	"github.com/ibuildthecloud/klum/pkg/controllers/user"
	"github.com/ibuildthecloud/klum/pkg/crd"
	"github.com/ibuildthecloud/klum/pkg/generated/controllers/klum.cattle.io"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/core"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/rbac"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/rancher/wrangler/pkg/start"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	Version    = "v0.0.0-dev"
	GitCommit  = "HEAD"
	cfg        user.Config
	kubeConfig string
)

func main() {
	app := cli.NewApp()
	app.Name = "klum"
	app.Version = fmt.Sprintf("%s (%s)", Version, GitCommit)
	app.Usage = "Kubernetes Lazy User Manager"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			EnvVar:      "KUBECONFIG",
			Destination: &kubeConfig,
		},
		cli.StringFlag{
			Name:        "namespace",
			EnvVar:      "NAMESPACE",
			Usage:       "Namespace to create secrets and SAs in",
			Value:       "klum",
			Destination: &cfg.Namespace,
		},
		cli.StringFlag{
			Name:        "context-name",
			Usage:       "Context name to put in Kubeconfigs",
			EnvVar:      "CONTEXT_NAME",
			Value:       "default",
			Destination: &cfg.ContextName,
		},
		cli.StringFlag{
			Name:        "server",
			Usage:       "The external server field to put in the Kubeconfigs",
			EnvVar:      "SERVER_NAME",
			Value:       "https://localhost:6443",
			Destination: &cfg.Server,
		},
		cli.StringFlag{
			Name:        "ca",
			Usage:       "The value of the CA data to put in the Kubeconfig",
			EnvVar:      "CA",
			Destination: &cfg.CA,
		},
		cli.StringFlag{
			Name:        "default-cluster-role",
			Usage:       "Default cluster-role to assign to users with no roles",
			EnvVar:      "DEFAULT_CLUSTER_ROLE",
			Value:       "cluster-admin",
			Destination: &cfg.DefaultClusterRole,
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	logrus.Info("Starting klum controller")
	ctx := signals.SetupSignalContext()

	restConfig, err := kubeconfig.GetNonInteractiveClientConfig(kubeConfig).ClientConfig()
	if err != nil {
		return err
	}

	if err := crd.Create(ctx, restConfig); err != nil {
		return err
	}

	core, err := core.NewFactoryFromConfig(restConfig)
	if err != nil {
		return err
	}

	klum, err := klum.NewFactoryFromConfigWithNamespace(restConfig, cfg.Namespace)
	if err != nil {
		return err
	}

	rbac, err := rbac.NewFactoryFromConfig(restConfig)
	if err != nil {
		return err
	}

	apply, err := apply.NewForConfig(restConfig)
	if err != nil {
		return nil
	}

	user.Register(ctx,
		cfg,
		apply,
		core.Core().V1().ServiceAccount(),
		rbac.Rbac().V1().ClusterRoleBinding(),
		rbac.Rbac().V1().RoleBinding(),
		core.Core().V1().Secret(),
		klum.Klum().V1alpha1().Kubeconfig(),
		klum.Klum().V1alpha1().User())

	if err := start.All(ctx, 2, klum, core, rbac); err != nil {
		logrus.Fatalf("Error starting: %s", err.Error())
	}

	<-ctx.Done()
	return nil
}
