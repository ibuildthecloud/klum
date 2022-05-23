module github.com/ibuildthecloud/klum

go 1.12

replace (
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009
	github.com/rancher/wrangler-api => github.com/dylanhitt/wrangler-api v0.7.0
)

require (
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/rancher/lasso v0.0.0-20210616224652-fc3ebd901c08
	github.com/rancher/wrangler v0.8.10
	github.com/rancher/wrangler-api v0.6.0
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.1
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver v0.18.0
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
)
