package rancher_macvlan

var pluginContainerImages = map[string]map[string]string{
	"Attaching Macvlan": {
		"kube-multus":                    "multus",
		"kube-macvlan-cni":               "staticMacvlan",
		"network-admission-deploy":       "admission",
		"network-controller":             "networkController",
		"kube-net-attach-def-controller": "nadController",
	},
	"Canal+Macvlan": {
		"cni-plugins":                    "hardenedCNIPlugins",
		"kube-multus":                    "multus",
		"kube-macvlan-cni":               "staticMacvlan",
		"network-admission-deploy":       "admission",
		"network-controller":             "networkController",
		"kube-net-attach-def-controller": "nadController",
		"install-cni":                    "calicoCNI",
		"flexvol-driver":                 "calicoFlexvol",
		"mount-bpffs":                    "calicoNode",
		"calico-node":                    "calicoNode",
		"kube-flannel":                   "flannel",
		"calico-kube-controllers":        "calicoControllers",
	},
	"Flannel+Macvlan": {
		"cni-plugins":                    "hardenedCNIPlugins",
		"kube-multus":                    "multus",
		"kube-macvlan-cni":               "staticMacvlan",
		"network-admission-deploy":       "admission",
		"network-controller":             "networkController",
		"kube-net-attach-def-controller": "nadController",
		"install-cni-plugin":             "flannelCNIPlugins",
		"install-cni":                    "flannel",
		"kube-flannel":                   "flannel",
	},
}

// Ensure all container images were tested
var testedContainerImages = map[string]map[string]bool{
	"Attaching Macvlan": {},
	"Canal+Macvlan":     {},
	"Flannel+Macvlan":   {},
}

var containerImageTestData = map[string]interface{}{
	"calicoCNI": map[string]string{
		"repository": "test/mirrored-calico-cni",
		"tag":        "v0.0.0",
	},
	"calicoFlexvol": map[string]string{
		"repository": "test/mirrored-calico-pod2daemon-flexvol",
		"tag":        "v0.0.0",
	},
	"calicoNode": map[string]string{
		"repository": "test/mirrored-calico-node",
		"tag":        "v0.0.0",
	},
	"calicoControllers": map[string]string{
		"repository": "test/mirrored-calico-kube-controllers",
		"tag":        "v0.0.0",
	},
	"flannel": map[string]string{
		"repository": "test/mirrored-flannelcni-flannel",
		"tag":        "v0.0.0",
	},
	"flannelCNIPlugins": map[string]string{
		"repository": "test/mirrored-flannelcni-flannel-cni-plugin",
		"tag":        "v0.0.0",
	},
	"hardenedCNIPlugins": map[string]string{
		"repository": "test/hardened-cni-plugins",
		"tag":        "v0.0.0",
	},
	"multus": map[string]string{
		"repository": "test/hardened-multus-cni",
		"tag":        "v0.0.0-rancher",
	},
	"networkController": map[string]string{
		"repository": "test/network-controller",
		"tag":        "v0.0.0",
	},
	"admission": map[string]string{
		"repository": "test/network-admission-deploy",
		"tag":        "v0.0.0",
	},
	"nadController": map[string]string{
		"repository": "test/k8s-net-attach-def-controller",
		"tag":        "v0.0.0",
	},
	"staticMacvlan": map[string]string{
		"repository": "test/static-macvlan-cni",
		"tag":        "v0.0.0",
	},
}
