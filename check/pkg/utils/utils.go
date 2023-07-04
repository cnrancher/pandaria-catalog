package utils

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
