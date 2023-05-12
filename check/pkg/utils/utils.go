package utils

import (
	"fmt"
	"reflect"
	"strings"

	"helm.sh/helm/v3/pkg/repo"
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

func CompareChart(v1 *repo.ChartVersion, v2 *repo.ChartVersion) error {
	if v1.Name != v2.Name {
		return fmt.Errorf("chart name %q != %q", v1.Name, v2.Name)
	}
	if !reflect.DeepEqual(v1.Annotations, v2.Annotations) {
		var a1, a2 []string
		for k, v := range v1.Annotations {
			a1 = append(a1, k+": "+v)
		}
		for k, v := range v2.Annotations {
			a2 = append(a2, k+": "+v)
		}
		return fmt.Errorf("chart %q annotations not equal: \n%v\n----\n%v",
			v1.Name, strings.Join(a1, "\n"), strings.Join(a2, "\n"))
	}
	if v1.AppVersion != v2.AppVersion {
		return fmt.Errorf("chart version %q != %q",
			v1.AppVersion, v2.AppVersion)
	}
	if v1.APIVersion != v2.APIVersion {
		return fmt.Errorf("chart APIVersion %q != %q",
			v1.APIVersion, v2.APIVersion)
	}
	if v1.Description != v2.Description {
		return fmt.Errorf("chart Description %q != %q",
			v1.Description, v2.Description)
	}
	if !reflect.DeepEqual(v1.Keywords, v2.Keywords) {
		return fmt.Errorf("chart keywords %v != %v",
			v1.Keywords, v2.Keywords)
	}

	return nil
}
