package rancher_flat_network

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rancher/hull/pkg/chart"
	"github.com/rancher/hull/pkg/checker"
	"github.com/rancher/hull/pkg/test"
	"github.com/rancher/hull/pkg/utils"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DefaultReleaseName = "rancher-flat-network"
	DefaultNamespace   = "cattle-flat-network"

	RancherNetworkControllerDeploymentExistsCheck = "RancherNetworkControllerDeploymentExistsCheck"
	FoundKey                                      = "found"
)

var ChartPath = utils.MustGetLatestChartVersionPathFromIndex("index.yaml", "rancher-flat-network", true)

var containerImages = map[string]string{
	// container name: values image name
	"init-multus-config":                      "deploy",
	"rancher-flat-network-multus":             "multus",
	"wait-webhook-tls-deploy":                 "deploy",
	"rancher-flat-network-operator-container": "flatNetworkOperator",
	"rancher-flat-network-tls-rollout":        "deploy",
	"rancher-flat-network-deploy":             "deploy",
	"flat-network-cni":                        "flatNetworkCNI",
}

// Ensure all container images were tested
var testedContainerImages = map[string]string{}

var suite = test.Suite{
	ChartPath: ChartPath,
	Cases: []test.Case{
		{
			Name:            "Using Defaults",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace),
		},
		{
			Name: "Set Values.arpPolicy",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("arpPolicy", "arping"),
		},
		{
			Name: "Set Values.proxyARP",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("proxyARP", true),
		},
		{
			Name: "Set Values.clusterCIDR",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("clusterCIDR", "10.200.0.0/16"),
		},
		{
			Name: "Set Values.serviceCIDR",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("serviceCIDR", "10.201.0.0/16"),
		},
		{
			Name: "Set Values.clusterType",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("clusterType", "K3s"),
		},
		{
			Name: "Set Values.flatNetworkOperator.replicas",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("flatNetworkOperator.replicas", "3"),
		},
		{
			Name: "Set Values.flatNetworkOperator.cattleResyncDefault",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("flatNetworkOperator.cattleResyncDefault", "120min"),
		},
		{
			Name: "Set Values.flatNetworkOperator.cattleDevMode",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("flatNetworkOperator.cattleDevMode", "true"),
		},
		{
			Name: "Set Values.flatNetworkOperator.limits.memory",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("flatNetworkOperator.limits.memory", "256M"),
		},
		{
			Name: "Set Values.flatNetworkOperator.limits.cpu",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("flatNetworkOperator.limits.cpu", "200m"),
		},
		{
			Name: "Set Values.deploy.rolloutSchedule",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("deploy.rolloutSchedule", "0 0 2 * *"),
		},
		{
			Name: "Set Values.multus.cni.version",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("multus.cniVersion", "0.3.1"),
		},
		{
			Name: "Set Values.multus.cni.version",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("multus.cni.version", "0.3.1"),
		},
		{
			Name: "Set Values.image & systemDefaultRegistry",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("flatNetworkOperator.image", map[string]string{
					"repository": "test/rancher-flat-network-operator",
					"tag":        "v0.0.0",
					"pullPolicy": "Always",
				}).Set("deploy.image",
				map[string]string{
					"repository": "test/rancher-flat-network-deploy",
					"tag":        "v0.0.0",
					"pullPolicy": "Always",
				}).Set("flatNetworkCNI.image",
				map[string]string{
					"repository": "test/rancher-flat-network-cni",
					"tag":        "v0.0.0",
					"pullPolicy": "Always",
				}).Set("multus.image", map[string]string{
				"repository": "test/hardened-multus-cni",
				"tag":        "v0.0.0",
				"pullPolicy": "Always",
			}).Set("global.cattle.systemDefaultRegistry", "harbor.test.io"),
		},
	},
	NamedChecks: []test.NamedCheck{
		{
			Name: "Check clusterCIDR & serviceCIDR",
			Covers: []string{
				".Values.clusterCIDR",
				".Values.serviceCIDR",
				".Values.arpPolicy",
				".Values.proxyARP",
			},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, deployment *appsv1.Deployment) {
					if deployment.Name != "rancher-flat-network-operator" {
						return
					}
					expectedClusterCIDR, ok := checker.RenderValue[string](tc, ".Values.clusterCIDR")
					if !ok {
						return
					}
					expectedServiceCIDR, ok := checker.RenderValue[string](tc, ".Values.serviceCIDR")
					if !ok {
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
						if c.Name != "rancher-flat-network-operator-container" {
							continue
						}
						for _, e := range c.Env {
							containersEnv[e.Name] = e.Value
						}
					}
					var set = map[string]string{
						"FLAT_NETWORK_CLUSTER_CIDR":   expectedClusterCIDR,
						"FLAT_NETWORK_SERVICE_CIDR":   expectedServiceCIDR,
						"FLAT_NETWORK_CNI_ARP_POLICY": expectedArpPolicy,
						"FLAT_NETWORK_CNI_PROXY_ARP":  expectedProxyARP,
					}
					for k, v := range set {
						assert.Equal(tc.T,
							v, containersEnv[k],
							"flat-network-operator doesn't have correct container env for key:%s, expected: %v, got: %v",
							k, v, containersEnv[k])
					}
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
					if obj.GetName() != "rancher-flat-network-multus-ds" {
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
						default:
							assert.Equal(tc.T,
								defaultExpected, s[1],
								"kube-multus container args key %q value expected %q, actual %q",
								key, defaultExpected, s[1])
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
									"rancher-flat-network-multus Volume cni HostPath Path expected %q, actual %q",
									MultusCniconfHostPathK3s, hp)
							} else {
								assert.Equal(tc.T,
									MultusCniconfHostPathDefault, hp,
									"rancher-flat-network-multus Volume cni HostPath Path expected %q, actual %q",
									MultusCniconfHostPathDefault, hp)
							}
							PodVolumeCniconfTested = true
						case "cnibin":
							if clusterType == "K3s" {
								assert.Equal(tc.T,
									MultusCnibinHostPathK3s, hp,
									"rancher-flat-network-multus Volume cnibin HostPath Path expected %q, actual %q",
									MultusCnibinHostPathK3s, hp)
							} else {
								assert.Equal(tc.T,
									MultusCnibinHostPathDefault, hp,
									"rancher-flat-network-multus Volume cnibin HostPath Path expected %q, actual %q",
									MultusCnibinHostPathDefault, hp)
							}
							PodVolumeCnibinTested = true
						}
					}
					if tc.T.Failed() {
						return
					}
					assert.Equal(tc.T, true,
						ContainerTested && PodVolumeCnibinTested && PodVolumeCniconfTested,
						"rancher-flat-network-multus-ds test not validated")
				}),
			},
		},
		{
			Name:   "Check multus.cni.version",
			Covers: []string{".Values.multus.cni.version"},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "rancher-flat-network-multus-ds" {
						return
					}
					expected, ok := checker.RenderValue[string](tc, ".Values.multus.cni.version")
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
						"rancher-flat-network-multus container env cni-version incorrect, expected:%s, got: %v",
						expected, actual)
				}),
			},
		},
		{
			Name: "Check flatNetworkOperator",
			Covers: []string{
				".Values.flatNetworkOperator.replicas",
				".Values.flatNetworkOperator.cattleResyncDefault",
				".Values.flatNetworkOperator.cattleDevMode",
				".Values.flatNetworkOperator.limits.memory",
				".Values.flatNetworkOperator.limits.cpu",
			},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, dep *appsv1.Deployment) {
					if dep.Name != "rancher-flat-network-operator" {
						return
					}
					expectedReplicas, _ := checker.RenderValue[string](tc, ".Values.flatNetworkOperator.replicas")
					expectedResyncDefault, _ := checker.RenderValue[string](tc, ".Values.flatNetworkOperator.cattleResyncDefault")
					expectedDevMode, _ := checker.RenderValue[string](tc, ".Values.flatNetworkOperator.cattleDevMode")
					expectedMemory, _ := checker.RenderValue[string](tc, ".Values.flatNetworkOperator.limits.memory")
					expectedCpu, _ := checker.RenderValue[string](tc, ".Values.flatNetworkOperator.limits.cpu")

					pod := dep.Spec.Template
					// set[expected][]any{KEY_NAME, actual}
					var set = map[any][]any{
						expectedReplicas:                   {"flatNetworkOperator.replicas", fmt.Sprintf("%v", *dep.Spec.Replicas)},
						expectedResyncDefault:              {"flatNetworkOperator.cattleResyncDefault", nil},
						expectedDevMode:                    {"flatNetworkOperator.cattleDevMode", nil},
						resource.MustParse(expectedMemory): {"flatNetworkOperator.limits.memory", *pod.Spec.Containers[0].Resources.Limits.Memory()},
						resource.MustParse(expectedCpu):    {"flatNetworkOperator.limits.cpu", *pod.Spec.Containers[0].Resources.Limits.Cpu()},
					}
					for _, e := range pod.Spec.Containers[0].Env {
						switch e.Name {
						case "CATTLE_RESYNC_DEFAULT":
							set[expectedResyncDefault][1] = e.Value
						case "CATTLE_DEV_MODE":
							set[expectedDevMode][1] = e.Value
						}
					}
					for e, a := range set {
						assert.Equal(tc.T, e, a[1],
							"rancher-flat-network-operator deployment [%v] incorrect, expected:%v, got: %v",
							a[0], e, a[1])
					}
				}),
			},
		},
		{
			Name: "Check image & systemDefaultRegistry",
			Covers: []string{
				".Values.flatNetworkOperator.image.repository",
				".Values.flatNetworkOperator.image.tag",
				".Values.flatNetworkOperator.image.pullPolicy",
				".Values.deploy.image.repository",
				".Values.deploy.image.tag",
				".Values.deploy.image.pullPolicy",
				".Values.flatNetworkCNI.image.repository",
				".Values.flatNetworkCNI.image.tag",
				".Values.flatNetworkCNI.image.pullPolicy",
				".Values.multus.image.repository",
				".Values.multus.image.tag",
				".Values.multus.image.pullPolicy",
				".Values.global.cattle.systemDefaultRegistry",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, pod corev1.PodTemplateSpec) {
					systemDefaultRegistry, ok := checker.RenderValue[string](tc, ".Values.global.cattle.systemDefaultRegistry")
					if !ok {
						tc.T.Logf("failed to get systemDefaultRegistry")
						return
					}
					testContainer := func(tc *checker.TestContext, container *corev1.Container) error {
						imageName, ok := containerImages[container.Name]
						if !ok {
							return fmt.Errorf("ignore unrecognized container: %v", container.Name)
						}
						repo, ok := checker.RenderValue[string](
							tc, fmt.Sprintf(".Values.%s.image.repository", imageName))
						if !ok {
							tc.T.Logf("failed to get image repo of container: %v", container.Name)
							return nil
						}
						tag, ok := checker.RenderValue[string](
							tc, fmt.Sprintf(".Values.%s.image.tag", imageName))
						if !ok {
							tc.T.Logf("failed to get image tag of container: %v", container.Name)
							return nil
						}
						pullPolicy, ok := checker.RenderValue[string](
							tc, fmt.Sprintf(".Values.%s.image.pullPolicy", imageName))
						if !ok {
							tc.T.Logf("failed to get image pullPolicy of container: %v", container.Name)
							return nil
						}
						expectedImage := fmt.Sprintf("%s/%s:%s", systemDefaultRegistry, repo, tag)
						actualImage := container.Image
						assert.Equal(tc.T, expectedImage, actualImage,
							"container image test failed, expected: %v, actual: %v",
							expectedImage, actualImage)
						assert.Equal(tc.T, pullPolicy, string(container.ImagePullPolicy),
							"container image pullPolicy test failed, expected: %v, actual: %v",
							pullPolicy, container.ImagePullPolicy)
						tc.T.Logf("expected: %v, actual: %v\n", expectedImage, actualImage)
						testedContainerImages[container.Name] = imageName
						return nil
					}

					for _, container := range pod.Spec.Containers {
						if err := testContainer(tc, &container); err != nil {
							tc.T.Error(err)
							return
						}
					}
					for _, container := range pod.Spec.InitContainers {
						if err := testContainer(tc, &container); err != nil {
							tc.T.Error(err)
							return
						}
					}
				}),
			},
		},
		{
			Name: "Check deploy.rolloutSchedule",
			Covers: []string{
				".Values.deploy.rolloutSchedule",
			},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, cronjob *batchv1.CronJob) {
					if cronjob.Name != "rancher-flat-network-tls-rollout" {
						return
					}
					expected, ok := checker.RenderValue[string](
						tc, ".Values.deploy.rolloutSchedule")
					if !ok {
						tc.T.Logf("failed to get rolloutSchedule of daemonset: %v", cronjob.Name)
						return
					}
					assert.Equal(tc.T, expected, cronjob.Spec.Schedule,
						"rancher-flat-network-tls-rollou schedule test failed, expected: %v, actual: %v",
						expected, cronjob.Spec.Schedule)
				}),
			},
		},
	},
}
