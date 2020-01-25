module github.com/ibuildthecloud/klum

go 1.12

replace github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009

require (
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/rancher/norman/v2 v2.0.0-20200111044641-76fd7a67396a
	github.com/rancher/wrangler v0.4.0
	github.com/rancher/wrangler-api v0.4.1
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.1
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	k8s.io/api v0.17.0
	k8s.io/apiextensions-apiserver v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
)
