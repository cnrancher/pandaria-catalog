package utils

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	RancherVersionAnnotationKey = "catalog.cattle.io/rancher-version"
	KubeVersionAnnotationKey    = "catalog.cattle.io/kube-version"
	HiddenAnnotationKey         = "catalog.cattle.io/hidden"

	NoRancherVersionFile             = "no-rancher-version.txt"
	NoKubeVersionFile                = "no-kube-version.txt"
	ImageCheckFailedFile             = "image-check-failed.txt"
	SystemDefaultRegistryCheckFailed = "system-default-registry-failed.txt"

	SEPARATOR = "===================="
)

func SaveSlice(name string, data []string) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(strings.Join(data, "\n"))
	if err != nil {
		logrus.Errorf("failed to write file: %v", err)
		return err
	}
	return nil
}
