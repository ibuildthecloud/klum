package main

import (
	"github.com/rancher/wrangler/pkg/cleanup"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	if err := cleanup.Cleanup("./pkg/apis"); err != nil {
		logrus.Fatal(err)
	}
	if err := os.RemoveAll("./pkg/generated"); err != nil {
		logrus.Fatal(err)
	}
}
