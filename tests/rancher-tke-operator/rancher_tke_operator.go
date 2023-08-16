package rancher_tke_operator

import (
	"github.com/cnrancher/pandaria-catalog/tests/common"
	"github.com/rancher/hull/pkg/chart"
	"github.com/rancher/hull/pkg/checker"
	"github.com/rancher/hull/pkg/test"
	"github.com/rancher/hull/pkg/utils"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var ChartPath = utils.MustGetLatestChartVersionPathFromIndex("index.yaml", "rancher-tke-operator", true)

const (
	DefaultReleaseName = "tke-config-operator"
	DefaultNamespace   = "cattle-system"

	RancherTkeOperatorDeployExistsCheck = "RancherTkeOperatorDeploymentExistsCheck"
	FoundKey                            = "found"
)

var suite = test.Suite{
	ChartPath: ChartPath,

	Cases: []test.Case{
		{
			Name: "Using Defaults",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace),
		},
		{
			Name: "Set .Values.tkeOperator.image.repository and .Values.tkeOperator.image.tag",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("tkeOperator", map[string]string{
					"image.repository": "test/tke-operator",
					"image.tag":        "v0.0.1-ent",
				}),
		},
		{
			Name: "Set .Values.httpProxy",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("httpProxy", testProxy),
		},
		{
			Name: "Set .Values.httpsProxy",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("httpsProxy", testProxy),
		},
		{
			Name: "Set .Values.noProxy",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("noProxy", testProxy),
		},
		{
			Name: "Set .Values.additionalTrustedCAs to true",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("additionalTrustedCAs", true),
		},
		{
			Name: "Set .Values.lockName",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("lockName", "test-lock-name"),
		},
		{
			Name: "Set .Values.lockNamespace",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("lockNamespace", "test-lock-namespace"),
		},
		{
			Name: "Set .Values.qps",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("qps", "test-qps"),
		},
		{
			Name: "Set .Values.burst",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("burst", "test-burst"),
		},
		{
			Name: "Set Values.global.systemDefaultRegistry",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				SetValue("global.systemDefaultRegistry", "test-registry"),
		},
	},

	NamedChecks: []test.NamedCheck{
		{
			Name:   "All Workloads Have Service Account",
			Checks: common.AllWorkloadsHaveServiceAccount,
		},
		{
			Name:   "All Workload Container Should Have SystemDefaultRegistryPrefix",
			Checks: common.AllContainerImagesShouldHaveSystemDefaultRegistryPrefix,
			Covers: []string{
				"Values.global.systemDefaultRegistry",
			},
		},
		{

			Name: "Check tke-operator image repository and tag",
			Covers: []string{
				".Values.tkeOperator.image.repository",
				".Values.tkeOperator.image.tag",
				".Values.global.systemDefaultRegistry",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					if obj.GetName() != "tke-config-operator" {
						return
					}

					checker.MapSet(tc, RancherTkeOperatorDeployExistsCheck, FoundKey, true)

					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}

					container := podTemplateSpec.Spec.Containers[0]

					systemDefaultRegistry := common.GetSystemDefaultRegistry(tc)

					rancherTkeOperatorRepo, _ := checker.RenderValue[string](tc, ".Values.tkeOperator.image.repository")
					rancherTkeOperatorTag, _ := checker.RenderValue[string](tc, ".Values.tkeOperator.image.tag")

					containerImage := rancherTkeOperatorRepo + ":" + rancherTkeOperatorTag

					expectedContainerImage := systemDefaultRegistry + containerImage

					assert.Equal(tc.T,
						expectedContainerImage, container.Image,
						"container %s of deployment %s does not have correct image: expected: %v got: %v",
						container.Name, obj.GetName(), expectedContainerImage, container.Image)

				}),
				rancherTkeOperatorDeployExistsCheck,
			},
		},
		{
			Name: "Check tke-operator args",
			Covers: []string{
				".Values.lockName",
				".Values.lockNamespace",
				".Values.qps",
				".Values.burst",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					if obj.GetName() != "tke-config-operator" {
						return
					}

					checker.MapSet(tc, RancherTkeOperatorDeployExistsCheck, FoundKey, true)

					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}

					container := podTemplateSpec.Spec.Containers[0]

					assertTkeOperatorArgsValues(tc, container.Args)

				}),
				rancherTkeOperatorDeployExistsCheck,
			},
		},
		{
			Name: "Check tke-operator env",
			Covers: []string{
				".Values.httpProxy",
				".Values.httpsProxy",
				".Values.noProxy",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					if obj.GetName() != "tke-config-operator" {
						return
					}

					checker.MapSet(tc, RancherTkeOperatorDeployExistsCheck, FoundKey, true)

					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}

					container := podTemplateSpec.Spec.Containers[0]

					assertTkeOperatorEnvValues(tc, container.Env)

				}),
				rancherTkeOperatorDeployExistsCheck,
			},
		},
		{
			Name: "Check additionalTrustedCAs",
			Covers: []string{
				".Values.additionalTrustedCAs",
			},

			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, deployment *appsv1.Deployment) {

					expectedAdditionalTrustedCAs, exists := checker.RenderValue[bool](tc, ".Values.additionalTrustedCAs")

					if exists {
						for _, container := range deployment.Spec.Template.Spec.Containers {
							assert.Equal(tc.T,
								expectedAdditionalTrustedCAs, true,
								"container %s of deployment %s does not have correct additionalTrustedCAs: expected: %v got: %v",
								container.Name, deployment.Name, expectedAdditionalTrustedCAs, true)
						}
					}
				}),
			},
		},
	},
}

var rancherTkeOperatorDeployExistsCheck = checker.Once(func(tc *checker.TestContext) {

	foundRancherTkeOperatorDeploy, _ := checker.MapGet[string, string, bool](tc, RancherTkeOperatorDeployExistsCheck, FoundKey)
	if !foundRancherTkeOperatorDeploy {
		tc.T.Error("err: tke-config-operator depoloyment not found")
	}
})

func assertTkeOperatorArgsValues(tc *checker.TestContext, args []string) {

	for _, arg := range args {
		expectedValue := ""
		argArray := strings.Split(arg, "=")

		if len(argArray) == 2 {
			switch argArray[0] {
			case "-lock_name":
				expectedValue, _ = checker.RenderValue[string](tc, ".Values.lockName")
			case "-lock_namespace":
				expectedValue, _ = checker.RenderValue[string](tc, ".Values.lockNamespace")
			case "-qps":
				expectedValue, _ = checker.RenderValue[string](tc, ".Values.qps")
			case "-burst":
				expectedValue, _ = checker.RenderValue[string](tc, ".Values.burst")
			}

			assert.Equal(tc.T,
				expectedValue, argArray[1],
				"container of tke-operator deployment does not have correct args value for arg:%s, expected: %v, got: %v",
				argArray[0], expectedValue, argArray[1])
		}
	}
}

func assertTkeOperatorEnvValues(tc *checker.TestContext, env []corev1.EnvVar) {

	for _, envVar := range env {
		expectedValue := ""

		switch envVar.Name {
		case "HTTP_PROXY":
			expectedValue, _ = checker.RenderValue[string](tc, ".Values.httpProxy")
		case "HTTPS_PROXY":
			expectedValue, _ = checker.RenderValue[string](tc, ".Values.httpsProxy")
		case "NO_PROXY":
			expectedValue, _ = checker.RenderValue[string](tc, ".Values.noProxy")
		default:
			expectedValue = envVar.Value
		}

		assert.Equal(tc.T,
			expectedValue, envVar.Value,
			"container of tk-operator deployment does not have correct env value for envKey:%s, expected: %v, got: %v",
			envVar.Name, expectedValue, envVar.Value)
	}
}
