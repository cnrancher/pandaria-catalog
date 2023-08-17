package rancher_macvlan

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/rancher/hull/pkg/chart"
	"github.com/rancher/hull/pkg/checker"
	"github.com/rancher/hull/pkg/test"
	"github.com/rancher/hull/pkg/utils"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DefaultReleaseName = "rancher-macvlan"
	DefaultNamespace   = "kube-system"
)

var ChartPath = utils.MustGetLatestChartVersionPathFromIndex("index.yaml", "rancher-macvlan", true)

var suite = test.Suite{
	ChartPath: ChartPath,
	Cases: []test.Case{
		{
			Name:            "Using Defaults",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace),
		},
		// Test podCIDR and iface for all plugins.
		{
			Name: "Set podCIDR & iface when plugin is 'Canal+Macvlan'",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("iface", "eth0").
				Set("podCIDR", "10.0.0.0/8").
				Set("plugin", "Canal+Macvlan"),
		},
		{
			Name: "Set podCIDR & iface when plugin is 'Flannel+Macvlan'",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("iface", "eth0").
				Set("podCIDR", "10.0.0.0/8").
				Set("plugin", "Flannel+Macvlan"),
		},
		// Test container images and systemDefaultRegistry for all plugins.
		{
			Name: "Set systemDefaultRegistry & image when plugin is 'AttachingMacvlan'",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("image", containerImageTestData).
				Set("global.cattle.systemDefaultRegistry", "docker.io").
				Set("plugin", "Attaching Macvlan"),
		},
		{
			Name: "Set systemDefaultRegistry & image when plugin is 'Canal+Macvlan'",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("image", containerImageTestData).
				Set("global.cattle.systemDefaultRegistry", "docker.io").
				Set("plugin", "Canal+Macvlan"),
		},
		{
			Name: "Set systemDefaultRegistry & image when plugin is 'Flannel+Macvlan'",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("image", containerImageTestData).
				Set("global.cattle.systemDefaultRegistry", "docker.io").
				Set("plugin", "Flannel+Macvlan"),
		},
		// Test privateRegistry.registryUrl & images for all plugins.
		{
			Name: "Set podCIDR, iface images and plugin to 'Attaching Macvlan'",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("privateRegistry.registryUrl", "example.io").
				Set("plugin", "Attaching Macvlan"),
		},
		{
			Name: "Set podCIDR, iface images and plugin to 'Canal+Macvlan'",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("privateRegistry.registryUrl", "example.io").
				Set("iface", "eth0"),
		},
		{
			Name: "Set podCIDR, iface images and plugin to 'Flannel+Macvlan'",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("privateRegistry.registryUrl", "example.io").
				Set("iface", "eth0"),
		},
		{
			Name: "Set Values.arpPolicy",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("arpPolicy", "arp_notify"),
		},
		{
			Name: "Set Values.proxyARP",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("proxyARP", true),
		},
		{
			Name: "Set Values.clusterType",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("clusterType", "K3s"),
		},
		{
			Name: "Set Values.multus.cniVersion",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("multus.cniVersion", "0.4.0"),
		},
		{
			Name: "Set Values.canal.vethmtu",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("canal.vethmtu", "1"),
		},
		{
			Name: "Set Values.ncResources.limits.memory",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("ncResources.limits.memory", "1024Mi"),
		},
		{
			Name: "Set Values.nadcResources.limits.memory",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("nadcResources.limits.memory", "1024Mi"),
		},
	},
	NamedChecks: []test.NamedCheck{
		{
			Name:   "Check iface",
			Covers: []string{".Values.iface"},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, configmap *corev1.ConfigMap) {
					if configmap.Name != "canal-config" {
						return
					}
					expected, ok := checker.RenderValue[string](tc, ".Values.iface")
					if !ok {
						return
					}
					actual := configmap.Data["canal_iface"]
					assert.Equal(tc.T, expected, actual, "iface expected: %v, actual: %v",
						expected, actual)
					tc.T.Logf("canal-config ConfigMap data iface expected: %v, actual: %v",
						expected, actual)
				}),
				checker.PerResource(func(tc *checker.TestContext, configmap *corev1.ConfigMap) {
					if configmap.Name != "kube-flannel-cfg" {
						return
					}
					expected, ok := checker.RenderValue[string](tc, ".Values.iface")
					if !ok {
						return
					}
					actual := configmap.Data["iface"]
					assert.Equal(tc.T, expected, actual, "iface expected: %v, actual: %v",
						expected, actual)
					tc.T.Logf("kube-flannel-cfg ConfigMap data iface expected: %v, actual: %v",
						expected, actual)
				}),
			},
		},
		{
			Name:   "Check podCIDR",
			Covers: []string{".Values.podCIDR"},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, configmap *corev1.ConfigMap) {
					if configmap.Name != "canal-config" {
						return
					}
					expected, ok := checker.RenderValue[string](tc, ".Values.podCIDR")
					if !ok {
						return
					}
					actualNetConf := configmap.Data["net-conf.json"]
					var d = map[string]any{}
					json.Unmarshal([]byte(actualNetConf), &d)
					actual, _ := d["Network"].(string)
					assert.Equal(tc.T, expected, actual, "podCIDR expected: %v, actual: %v",
						expected, actual)
					tc.T.Logf("canal-config ConfigMap data podCIDR expected: %v, actual: %v",
						expected, actual)
				}),
				checker.PerResource(func(tc *checker.TestContext, configmap *corev1.ConfigMap) {
					if configmap.Name != "kube-flannel-cfg" {
						return
					}
					expected, ok := checker.RenderValue[string](tc, ".Values.podCIDR")
					if !ok {
						return
					}
					actualNetConf := configmap.Data["net-conf.json"]
					var d = map[string]any{}
					json.Unmarshal([]byte(actualNetConf), &d)
					actual, _ := d["Network"].(string)
					assert.Equal(tc.T, expected, actual, "podCIDR expected: %v, actual: %v",
						expected, actual)
					tc.T.Logf("kube-flannel-cfg ConfigMap data podCIDR expected: %v, actual: %v",
						expected, actual)
				}),
			},
		},
		{
			Name: "Check arpPolicy and proxyARP",
			Covers: []string{
				".Values.arpPolicy",
				".Values.proxyARP",
			},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, deployment *appsv1.Deployment) {
					if deployment.Name != "network-controller" {
						return
					}
					expectedArpPolicy, ok := checker.RenderValue[string](tc, ".Values.arpPolicy")
					if !ok {
						return
					}
					expectedProxyARPBool, _ := checker.RenderValue[bool](tc, ".Values.proxyARP")
					expectedProxyARP := strconv.FormatBool(expectedProxyARPBool)
					containersEnv := map[string]string{}
					for _, c := range deployment.Spec.Template.Spec.Containers {
						if c.Name != "network-controller" {
							continue
						}
						for _, e := range c.Env {
							containersEnv[e.Name] = e.Value
						}
					}
					var (
						arpPolicyEnvKey = "PANDARIA_MACVLAN_CNI_ARP_POLICY"
						proxyARPEnvKey  = "PANDARIA_MACVLAN_CNI_PROXY_ARP"
					)
					assert.Equal(tc.T,
						expectedArpPolicy, containersEnv[arpPolicyEnvKey],
						"network-controller doesn't have correct container env for key:%s, expected: %v, got: %v",
						arpPolicyEnvKey, expectedArpPolicy, containersEnv[arpPolicyEnvKey])
					assert.Equal(tc.T,
						expectedProxyARP, containersEnv[proxyARPEnvKey],
						"network-controller doesn't have correct container env for key:%s, expected: %v, got: %v",
						proxyARPEnvKey, expectedProxyARP, containersEnv[proxyARPEnvKey])
				}),
			},
		},
		{
			Name:   "Check clusterType",
			Covers: []string{".Values.clusterType"},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					clusterType, ok := checker.RenderValue[string](tc, ".Values.clusterType")
					if !ok {
						return
					}
					if obj.GetName() != "kube-multus-ds" {
						return
					}
					var (
						key                    = "--multus-kubeconfig-file-host"
						k3sExpected            = "/var/lib/rancher/k3s/agent/etc/cni/net.d/multus.d/multus.kubeconfig"
						defaultExpected        = "/etc/cni/net.d/multus.d/multus.kubeconfig"
						ContainerTested        = false
						PodVolumeCniconfTested = false
						PodVolumeCnibinTested  = false
					)
					for _, a := range podTemplateSpec.Spec.Containers[0].Args {
						s := strings.Split(a, "=")
						if len(s) != 2 {
							continue
						}
						if s[0] != key {
							continue
						}
						switch clusterType {
						case "K3s":
							assert.Equal(tc.T,
								k3sExpected, s[1],
								"kube-multus container args key %q value expected %q, actual %q",
								key, k3sExpected, s[1])
							tc.T.Logf("expected: %v, actual %v", k3sExpected, s[1])
						default:
							assert.Equal(tc.T,
								defaultExpected, s[1],
								"kube-multus container args key %q value expected %q, actual %q",
								key, defaultExpected, s[1])
							tc.T.Logf("expected: %v, actual %v", defaultExpected, s[1])
						}
						ContainerTested = true
					}
					if tc.T.Failed() {
						return
					}

					var (
						MultusCniconfHostPathK3s     = "/var/lib/rancher/k3s/agent/etc/cni/net.d"
						MultusCniconfHostPathDefault = "/etc/cni/net.d"
						MultusCnibinHostPathK3s      = "/var/lib/rancher/k3s/data/current/bin"
						MultusCnibinHostPathDefault  = "/opt/cni/bin"
					)
					for _, v := range podTemplateSpec.Spec.Volumes {
						if v.HostPath == nil {
							continue
						}
						hp := v.HostPath.Path
						switch v.Name {
						case "cni":
							if clusterType == "K3s" {
								assert.Equal(tc.T,
									MultusCniconfHostPathK3s, hp,
									"kube-multus-ds Volume cni HostPath Path expected %q, actual %q",
									MultusCniconfHostPathK3s, hp)
								tc.T.Logf("expected: %v, actual %v", MultusCniconfHostPathK3s, hp)
							} else {
								assert.Equal(tc.T,
									MultusCniconfHostPathDefault, hp,
									"kube-multus-ds Volume cni HostPath Path expected %q, actual %q",
									MultusCniconfHostPathDefault, hp)
								tc.T.Logf("expected: %v, actual %v", MultusCniconfHostPathDefault, hp)
							}
							PodVolumeCniconfTested = true
						case "cnibin":
							if clusterType == "K3s" {
								assert.Equal(tc.T,
									MultusCnibinHostPathK3s, hp,
									"kube-multus-ds Volume cnibin HostPath Path expected %q, actual %q",
									MultusCnibinHostPathK3s, hp)
								tc.T.Logf("expected: %v, actual %v", MultusCnibinHostPathK3s, hp)
							} else {
								assert.Equal(tc.T,
									MultusCnibinHostPathDefault, hp,
									"kube-multus-ds Volume cnibin HostPath Path expected %q, actual %q",
									MultusCnibinHostPathDefault, hp)
								tc.T.Logf("expected: %v, actual %v", MultusCnibinHostPathDefault, hp)
							}
							PodVolumeCnibinTested = true
						}
					}
					if tc.T.Failed() {
						return
					}
					assert.Equal(tc.T, true,
						ContainerTested && PodVolumeCnibinTested && PodVolumeCniconfTested,
						"kube-multus-ds test not validated")
				}),
			},
		},
		{
			Name:   "Check multus.cniVersion",
			Covers: []string{".Values.multus.cniVersion"},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "kube-multus-ds" {
						return
					}
					expected, ok := checker.RenderValue[string](tc, ".Values.multus.cniVersion")
					if !ok {
						return
					}
					var actual string
					for _, arg := range podTemplateSpec.Spec.Containers[0].Args {
						if strings.HasPrefix(arg, "--cni-version=") {
							a := strings.Split(arg, "=")
							if len(a) != 2 {
								continue
							}
							actual = a[1]
						}
					}
					assert.Equal(tc.T, expected, actual,
						"kube-multus container env cni-version incorrect, expected:%s, got: %v",
						expected, actual)
				}),
			},
		},
		{
			Name:   "Check canal.vethmtu",
			Covers: []string{".Values.canal.vethmtu"},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, configmap *corev1.ConfigMap) {
					if configmap.Name != "canal-config" {
						return
					}
					expected, ok := checker.RenderValue[string](tc, ".Values.canal.vethmtu")
					if !ok {
						tc.T.Error("failed to get .Values.canal.vethmtu")
						return
					}
					actual := configmap.Data["veth_mtu"]
					assert.Equal(tc.T, expected, actual, "canal.vethmtu expected: %v, actual: %v",
						expected, actual)
				}),
			},
		},
		{
			Name:   "Check ncResources.limits.memory",
			Covers: []string{".Values.ncResources.limits.memory"},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "network-controller" {
						return
					}
					expectedVal, ok := checker.RenderValue[string](tc, ".Values.ncResources.limits.memory")
					if !ok {
						return
					}
					expected := resource.MustParse(expectedVal)
					actual := *podTemplateSpec.Spec.Containers[0].Resources.Limits.Memory()
					assert.True(tc.T, actual.Equal(expected),
						"network-controller container resource limits memory incorrect, expected:%s, got: %v",
						expected, actual)
				}),
			},
		},
		{
			Name:   "Check nadcResources.limits.memory",
			Covers: []string{".Values.nadcResources.limits.memory"},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "kube-net-attach-def-controller" {
						return
					}
					expectedVal, ok := checker.RenderValue[string](tc, ".Values.nadcResources.limits.memory")
					if !ok {
						return
					}
					expected := resource.MustParse(expectedVal)
					actual := *podTemplateSpec.Spec.Containers[0].Resources.Limits.Memory()
					assert.True(tc.T, actual.Equal(expected),
						"kube-net-attach-def-controller container resource limits memory incorrect, expected:%s, got: %v",
						expected, actual)
				}),
			},
		},
		{
			Name: "Check image & systemDefaultRegistry & plugin",
			Covers: []string{
				".Values.plugin",
				".Values.image.calicoCNI.repository",
				".Values.image.calicoCNI.tag",
				".Values.image.calicoFlexvol.repository",
				".Values.image.calicoFlexvol.tag",
				".Values.image.calicoNode.repository",
				".Values.image.calicoNode.tag",
				".Values.image.calicoControllers.repository",
				".Values.image.calicoControllers.tag",
				".Values.image.flannel.repository",
				".Values.image.flannel.tag",
				".Values.image.flannelCNIPlugins.repository",
				".Values.image.flannelCNIPlugins.tag",
				".Values.image.hardenedCNIPlugins.repository",
				".Values.image.hardenedCNIPlugins.tag",
				".Values.image.multus.repository",
				".Values.image.multus.tag",
				".Values.image.admission.repository",
				".Values.image.admission.tag",
				".Values.image.nadController.repository",
				".Values.image.nadController.tag",
				".Values.image.networkController.repository",
				".Values.image.networkController.tag",
				".Values.image.staticMacvlan.repository",
				".Values.image.staticMacvlan.tag",
				".Values.global.cattle.systemDefaultRegistry",
			},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, deployment *appsv1.Deployment) {
					systemDefaultRegistry, ok := checker.RenderValue[string](tc, ".Values.global.cattle.systemDefaultRegistry")
					if !ok {
						return
					}
					plugin, ok := checker.RenderValue[string](tc, ".Values.plugin")
					if !ok {
						return
					}
					tc.T.Logf("deployment: %v, plugin: %v", deployment.GetName(), plugin)
					for _, container := range deployment.Spec.Template.Spec.Containers {
						testContainerImages(tc, plugin, systemDefaultRegistry, &container)
					}
					for _, container := range deployment.Spec.Template.Spec.InitContainers {
						testContainerImages(tc, plugin, systemDefaultRegistry, &container)
					}
				}),
				checker.PerResource(func(tc *checker.TestContext, daemonSet *appsv1.DaemonSet) {
					systemDefaultRegistry, ok := checker.RenderValue[string](tc, ".Values.global.cattle.systemDefaultRegistry")
					if !ok {
						return
					}
					plugin, ok := checker.RenderValue[string](tc, ".Values.plugin")
					if !ok {
						return
					}
					tc.T.Logf("daemonSet: %v", daemonSet.GetName())
					for _, container := range daemonSet.Spec.Template.Spec.Containers {
						testContainerImages(tc, plugin, systemDefaultRegistry, &container)
					}
					for _, container := range daemonSet.Spec.Template.Spec.InitContainers {
						testContainerImages(tc, plugin, systemDefaultRegistry, &container)
					}
				}),
			},
		},
		{
			Name:   "Check privateRegistry.registryUrl",
			Covers: []string{".Values.privateRegistry.registryUrl"},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, deployment *appsv1.Deployment) {
					registryUrl, ok := checker.RenderValue[string](tc, ".Values.privateRegistry.registryUrl")
					if !ok {
						return
					}
					plugin, ok := checker.RenderValue[string](tc, ".Values.plugin")
					if !ok {
						return
					}
					tc.T.Logf("deployment: %v, plugin: %v", deployment.GetName(), plugin)
					for _, container := range deployment.Spec.Template.Spec.Containers {
						testContainerImages(tc, plugin, registryUrl, &container)
					}
					for _, container := range deployment.Spec.Template.Spec.InitContainers {
						testContainerImages(tc, plugin, registryUrl, &container)
					}
				}),
				checker.PerResource(func(tc *checker.TestContext, daemonSet *appsv1.Deployment) {
					registryUrl, ok := checker.RenderValue[string](tc, ".Values.privateRegistry.registryUrl")
					if !ok {
						return
					}
					plugin, ok := checker.RenderValue[string](tc, ".Values.plugin")
					if !ok {
						return
					}
					tc.T.Logf("daemonSet: %v", daemonSet.GetName())
					for _, container := range daemonSet.Spec.Template.Spec.Containers {
						testContainerImages(tc, plugin, registryUrl, &container)
					}
					for _, container := range daemonSet.Spec.Template.Spec.InitContainers {
						testContainerImages(tc, plugin, registryUrl, &container)
					}
				}),
			},
		},
	},
}

func printObject(a any) string {
	b, _ := json.MarshalIndent(a, "", "  ")
	return string(b)
}

func testContainerImages(tc *checker.TestContext, plugin, registry string, container *corev1.Container) {
	imageName, ok := pluginContainerImages[plugin][container.Name]
	if !ok {
		tc.T.Errorf("ignore unrecognized container: %v", container.Name)
		return
	}
	repo, ok := checker.RenderValue[string](
		tc, fmt.Sprintf(".Values.image.%s.repository", imageName))
	if !ok {
		tc.T.Errorf("failed to get image repo of container: %v", container.Name)
		return
	}
	tag, ok := checker.RenderValue[string](
		tc, fmt.Sprintf(".Values.image.%s.tag", imageName))
	if !ok {
		tc.T.Errorf("failed to get image tag of container: %v", container.Name)
		return
	}
	expectedImage := fmt.Sprintf("%s/%s:%s", registry, repo, tag)
	actualImage := container.Image
	assert.Equal(tc.T, expectedImage, actualImage,
		"container image test failed, expected: %v, actual: %v",
		expectedImage, actualImage)
	tc.T.Logf("expected: %v, actual: %v", expectedImage, actualImage)
	testedContainerImages[plugin][container.Name] = true
}
