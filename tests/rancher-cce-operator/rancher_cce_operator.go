package rancher_cce_operator

import (
	"strings"

	"github.com/cnrancher/pandaria-catalog/tests/common"
	"github.com/rancher/hull/pkg/chart"
	"github.com/rancher/hull/pkg/checker"
	"github.com/rancher/hull/pkg/test"
	"github.com/rancher/hull/pkg/utils"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ChartPath = utils.MustGetLatestChartVersionPathFromIndex("index.yaml", "rancher-cce-operator", true)

const (
	DefaultReleaseName = "cce-config-operator"
	DefaultNamespace   = "cattle-system"

	RancherCceOperatorDeployExistsCheck = "RancherCceOperatorDeploymentExistsCheck"
	FoundKey                            = "found"
)

var suite = test.Suite{
	ChartPath: ChartPath,

	Cases: []test.Case{
		{
			Name:            "Using Defaults",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace),
		},
		{
			Name: "Set .Values.cceOperator.image.repository and .Values.cceOperator.image.tag",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("cceOperator", map[string]string{
					"image.repository": "test/cce-operator",
					"image.tag":        "v0.0.1",
				}),
		},
		{
			Name: "Set .Values.httpProxy",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("httpProxy", "http://127.0.0.1"),
		},
		{
			Name: "Set .Values.httpsProxy",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("httpsProxy", "http://127.0.0.1"),
		},
		{
			Name: "Set .Values.noProxy",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("noProxy", "127.0.0.1"),
		},
		{
			Name: "Set .Values.additionalTrustedCAs to true",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("additionalTrustedCAs", true),
		},
		{
			Name: "Set .Values.nodeSelector",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("nodeSelector", map[string]interface{}{
					"diskType": "ssd",
				}),
		},
		{
			Name: "Set .Values.tolerations",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("tolerations", []corev1.Toleration{
					{
						Key:      "example-key",
						Operator: corev1.TolerationOpExists,
						Value:    "test",
						Effect:   corev1.TaintEffectNoSchedule,
					},
				}),
		},
		{
			Name: "Set .Values.priorityClassName",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("priorityClassName", "high-priority"),
		},
		{
			Name: "Set Values.global.cattle.systemDefaultRegistry",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				SetValue("global.cattle.systemDefaultRegistry", "test-registry"),
		},
	},

	NamedChecks: []test.NamedCheck{
		{
			Name:   "All Workloads Have Service Account",
			Checks: common.AllWorkloadsHaveServiceAccount,
		},
		{

			Name: "Check cce-operator image repository and tag",
			Covers: []string{
				".Values.cceOperator.image.repository",
				".Values.cceOperator.image.tag",
				".Values.global.cattle.systemDefaultRegistry",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					if obj.GetName() != "cce-config-operator" {
						return
					}

					checker.MapSet(tc, RancherCceOperatorDeployExistsCheck, FoundKey, true)

					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}

					container := podTemplateSpec.Spec.Containers[0]

					systemDefaultRegistry, _ := checker.RenderValue[string](tc, ".Values.global.cattle.systemDefaultRegistry")
					rancherCceOperatorRepo, _ := checker.RenderValue[string](tc, ".Values.cceOperator.image.repository")
					rancherCceOperatorTag, _ := checker.RenderValue[string](tc, ".Values.cceOperator.image.tag")
					if systemDefaultRegistry != "" {
						systemDefaultRegistry += "/"
					}
					containerImage := rancherCceOperatorRepo + ":" + rancherCceOperatorTag
					expectedContainerImage := systemDefaultRegistry + containerImage
					assert.Equal(tc.T,
						expectedContainerImage, container.Image,
						"container %s of deployment %s does not have correct image: expected: %v got: %v",
						container.Name, obj.GetName(), expectedContainerImage, container.Image)
				}),
				rancherCceOperatorDeployExistsCheck,
			},
		},
		{
			Name: "Check cce-operator args",
			Covers: []string{
				".Values.lockName",
				".Values.lockNamespace",
				".Values.qps",
				".Values.burst",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					if obj.GetName() != "cce-config-operator" {
						return
					}

					checker.MapSet(tc, RancherCceOperatorDeployExistsCheck, FoundKey, true)

					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}

					container := podTemplateSpec.Spec.Containers[0]

					assertCceOperatorArgsValues(tc, container.Args)

				}),
				rancherCceOperatorDeployExistsCheck,
			},
		},
		{
			Name: "Check cce-operator env",
			Covers: []string{
				".Values.httpProxy",
				".Values.httpsProxy",
				".Values.noProxy",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					if obj.GetName() != "cce-config-operator" {
						return
					}

					checker.MapSet(tc, RancherCceOperatorDeployExistsCheck, FoundKey, true)

					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}

					container := podTemplateSpec.Spec.Containers[0]

					assertCceOperatorEnvValues(tc, container.Env)

				}),
				rancherCceOperatorDeployExistsCheck,
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
		{
			Name: "Check nodeSelector",
			Covers: []string{
				".Values.nodeSelector",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "cce-config-operator" {
						return
					}

					nodeSelector := podTemplateSpec.Spec.NodeSelector
					if len(nodeSelector) == 0 {
						return
					}
					diskType, ok := nodeSelector["diskType"]
					if !ok {
						return
					}
					assert.Equal(tc.T, "ssd", diskType)
				}),
			},
		},
		{
			Name: "Check tolerations",
			Covers: []string{
				".Values.tolerations",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "cce-config-operator" {
						return
					}

					tolerations := podTemplateSpec.Spec.Tolerations
					assert.NotNil(tc.T, tolerations)
					assert.True(tc.T, len(tolerations) > 0)
					for _, t := range tolerations {
						if t.Key != "example-key" {
							continue
						}
						assert.Equal(tc.T, "Exists", string(t.Operator))
						assert.Equal(tc.T, "NoSchedule", string(t.Effect))
					}
				}),
			},
		},
		{
			Name: "Check priorityClassName",
			Covers: []string{
				".Values.priorityClassName",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "cce-config-operator" {
						return
					}
					priorityClassName := podTemplateSpec.Spec.PriorityClassName
					if priorityClassName == "" {
						return
					}
					assert.Equal(tc.T, "high-priority", priorityClassName)
				}),
			},
		},
	},
}

var rancherCceOperatorDeployExistsCheck = checker.Once(func(tc *checker.TestContext) {
	foundRancherCceOperatorDeploy, _ := checker.MapGet[string, string, bool](tc, RancherCceOperatorDeployExistsCheck, FoundKey)
	if !foundRancherCceOperatorDeploy {
		tc.T.Error("err: cce-config-operator depoloyment not found")
	}
})

func assertCceOperatorArgsValues(tc *checker.TestContext, args []string) {

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
				"container of cce-operator deployment does not have correct args value for arg:%s, expected: %v, got: %v",
				argArray[0], expectedValue, argArray[1])
		}
	}
}

func assertCceOperatorEnvValues(tc *checker.TestContext, env []corev1.EnvVar) {
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
			"container of cce-operator deployment does not have correct env value for envKey:%s, expected: %v, got: %v",
			envVar.Name, expectedValue, envVar.Value)
	}
}
