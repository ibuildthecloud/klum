package main

import (
	"os"

	"github.com/jadolg/klum/pkg/apis/klum.cattle.io/v1alpha1"
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/jadolg/klum/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"klum.cattle.io": {
				Types: []interface{}{
					v1alpha1.User{},
					v1alpha1.Kubeconfig{},
					v1alpha1.UserSyncGithub{},
				},
				GenerateTypes: true,
			},
		},
	})
}
