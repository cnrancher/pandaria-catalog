package rancher_ack_operator

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

var ChartPath = utils.MustGetLatestChartVersionPathFromIndex("index.yaml", "rancher-ack-operator", true)

const (
	DefaultReleaseName = "ack-config-operator"
	DefaultNamespace   = "cattle-system"

	RancherAckOperatorDeployExistsCheck = "RancherAckOperatorDeploymentExistsCheck"
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
			Name: "Set .Values.ackOperator.image.repository and .Values.ackOperator.image.tag",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("ackOperator", map[string]string{
					"image.repository": "test/ack-operator",
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
			Name: "Set .Values.leaderElect to true",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("leaderElect", true),
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
			Covers: []string{
				"Values.global.systemDefaultRegistry",
			},
		},
		{
			Name:   "All Workload Container Should Have SystemDefaultRegistryPrefix",
			Checks: common.AllContainerImagesShouldHaveSystemDefaultRegistryPrefix,
			Covers: []string{
				"Values.global.systemDefaultRegistry",
			},
		},
		{
			Name: "Check ack-operator image repository and tag",
			Covers: []string{
				".Values.ackOperator.image.repository",
				".Values.ackOperator.image.tag",
				".Values.global.systemDefaultRegistry",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "ack-config-operator" {
						return
					}
					checker.MapSet(tc, RancherAckOperatorDeployExistsCheck, FoundKey, true)
					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))
					if !ok {
						return
					}
					container := podTemplateSpec.Spec.Containers[0]
					systemDefaultRegistry := common.GetSystemDefaultRegistry(tc)
					rancherAckOperatorRepo, _ := checker.RenderValue[string](tc, ".Values.ackOperator.image.repository")
					rancherAckOperatorTag, _ := checker.RenderValue[string](tc, ".Values.ackOperator.image.tag")
					containerImage := rancherAckOperatorRepo + ":" + rancherAckOperatorTag
					expectedContainerImage := systemDefaultRegistry + containerImage
					assert.Equal(tc.T,
						expectedContainerImage, container.Image,
						"container %s of deployment %s does not have correct image: expected: %v got: %v",
						container.Name, obj.GetName(), expectedContainerImage, container.Image)

				}),
				rancherAckOperatorDeployExistsCheck,
			},
		},
		{
			Name: "Check ack-operator args",
			Covers: []string{
				".Values.leaderElect",
				".Values.lockName",
				".Values.lockNamespace",
				".Values.qps",
				".Values.burst",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "ack-config-operator" {
						return
					}
					checker.MapSet(tc, RancherAckOperatorDeployExistsCheck, FoundKey, true)
					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}
					container := podTemplateSpec.Spec.Containers[0]
					assertAckOperatorArgsValues(tc, container.Args)
				}),
				rancherAckOperatorDeployExistsCheck,
			},
		},
		{
			Name: "Check ack-operator env",
			Covers: []string{
				".Values.httpProxy",
				".Values.httpsProxy",
				".Values.noProxy",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "ack-config-operator" {
						return
					}
					checker.MapSet(tc, RancherAckOperatorDeployExistsCheck, FoundKey, true)
					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}
					container := podTemplateSpec.Spec.Containers[0]
					assertAckOperatorEnvValues(tc, container.Env)

				}),
				rancherAckOperatorDeployExistsCheck,
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

var rancherAckOperatorDeployExistsCheck = checker.Once(func(tc *checker.TestContext) {
	foundRancherAckOperatorDeploy, _ := checker.MapGet[string, string, bool](tc, RancherAckOperatorDeployExistsCheck, FoundKey)
	if !foundRancherAckOperatorDeploy {
		tc.T.Error("err: ack-config-operator depoloyment not found")
	}
})

func assertAckOperatorArgsValues(tc *checker.TestContext, args []string) {
	for _, arg := range args {
		expectedValue := ""
		argArray := strings.Split(arg, "=")
		if len(argArray) == 2 {
			switch argArray[0] {
			case "-leader_elect":
				expectedValue, _ = checker.RenderValue[string](tc, ".Values.leaderElect")
				if expectedValue == "" {
					expectedValue = "true"
				}
			case "-lock_name":
				expectedValue, _ = checker.RenderValue[string](tc, ".Values.lockName")
			case "-lock_namespace":
				expectedValue, _ = checker.RenderValue[string](tc, ".Values.lockNamespace")
				if expectedValue == "" {
					expectedValue = "cattle-system"
				}
			case "-qps":
				expectedValue, _ = checker.RenderValue[string](tc, ".Values.qps")
			case "-burst":
				expectedValue, _ = checker.RenderValue[string](tc, ".Values.burst")
			}
			assert.Equal(tc.T,
				expectedValue, argArray[1],
				"container of ack-operator deployment does not have correct args value for arg:%s, expected: %v, got: %v",
				argArray[0], expectedValue, argArray[1])
		}
	}
}

func assertAckOperatorEnvValues(tc *checker.TestContext, env []corev1.EnvVar) {
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
			"container of ack-operator deployment does not have correct env value for envKey:%s, expected: %v, got: %v",
			envVar.Name, expectedValue, envVar.Value)
	}
}
