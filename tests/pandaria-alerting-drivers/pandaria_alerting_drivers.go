package pandaria_alerting_drivers

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
)

var ChartPath = utils.MustGetLatestChartVersionPathFromIndex("index.yaml", "pandaria-alerting-drivers", true)

const (
	DefaultReleaseName = "pandaria-alerting-drivers"
	DefaultNamespace   = "cattle-monitoring-system"

	PandariaAlertingDriversDeployExistsCheck = "PandariaAlertingDriversDeploymentExistsCheck"
	FoundKey                                 = "found"
)

var suite = test.Suite{
	ChartPath: ChartPath,

	Cases: []test.Case{
		{
			Name: "Using Defaults",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace),
		},
		{
			Name: "Set .Values.enabled to true",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("enabled", true),
		},
		{
			Name: "Set .Values.image.repository and .Values.image.tag",
			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("image", map[string]string{
					"repository": "test/webhook-receiver",
					"tag":        "v0.0.1-ent",
				}),
		},
		{
			Name: "Set .Values.replicas to 2",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("replicas", 2),
		},
		{
			Name: "Set Values.resources",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("resources", testResources),
		},
		{
			Name: "Set Values.nodeSelector",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("nodeSelector", testNodeSelector),
		},
		{
			Name: "Set Values.tolerations",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("tolerations", testTolerations),
		},
		{
			Name: "Set Values.affinity",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("affinity", testAffinity),
		},
		{
			Name: "Set .Values.port to 9090",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("port", 9090),
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
			Name:   "All Workloads Have Node Selectors And Tolerations For OS",
			Checks: common.AllWorkloadsHaveNodeSelectorsAndTolerationsForOS,
		},
		{
			Name:   "All Workload Container Should Have SystemDefaultRegistryPrefix",
			Checks: common.AllContainerImagesShouldHaveSystemDefaultRegistryPrefix,
			Covers: []string{
				"Values.global.systemDefaultRegistry",
			},
		},
		{
			Name: "Check All Workloads Have NodeSelector As Per Given Value",
			Covers: []string{
				".Values.nodeSelector",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					nodeSelectorAddedByValues, _ := checker.RenderValue[map[string]string](tc, ".Values.nodeSelector")

					expectedNodeSelector := map[string]string{}

					for k, v := range nodeSelectorAddedByValues {
						expectedNodeSelector[k] = v
					}

					for k, v := range defaultNodeSelector {
						expectedNodeSelector[k] = v
					}

					assert.Equal(tc.T,
						expectedNodeSelector, podTemplateSpec.Spec.NodeSelector,
						"workload %s (type: %T) does not have correct nodeSelectors, expected: %v got: %v",
						obj.GetName(), obj, expectedNodeSelector, podTemplateSpec.Spec.NodeSelector,
					)
				}),
			},
		},
		{
			Name: "Check All Workloads Have Tolerations As Per Given Value",
			Covers: []string{
				".Values.tolerations",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					tolerationsAddedByValues, _ := checker.RenderValue[[]corev1.Toleration](tc, ".Values.tolerations")

					expectedTolerations := append(defaultTolerations, tolerationsAddedByValues...)
					if len(expectedTolerations) == 0 {
						expectedTolerations = nil
					}

					assert.Equal(tc.T,
						expectedTolerations, podTemplateSpec.Spec.Tolerations,
						"workload %s (type: %T) does not have correct tolerations, expected: %v got: %v",
						obj.GetName(), obj, expectedTolerations, podTemplateSpec.Spec.Tolerations,
					)
				}),
			},
		},
		{

			Name: "Check webhook-receiver image repository and tag",
			Covers: []string{
				".Values.image.repository",
				".Values.image.tag",
				".Values.global.systemDefaultRegistry",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					if obj.GetName() != "pandaria-alerting-drivers" {
						return
					}

					checker.MapSet(tc, PandariaAlertingDriversDeployExistsCheck, FoundKey, true)

					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}

					container := podTemplateSpec.Spec.Containers[0]

					systemDefaultRegistry := common.GetSystemDefaultRegistry(tc)

					pandariaAlertingDriversRepo, _ := checker.RenderValue[string](tc, ".Values.image.repository")
					pandariaAlertingDriversTag, _ := checker.RenderValue[string](tc, ".Values.image.tag")

					containerImage := pandariaAlertingDriversRepo + ":" + pandariaAlertingDriversTag

					expectedContainerImage := systemDefaultRegistry + containerImage

					assert.Equal(tc.T,
						expectedContainerImage, container.Image,
						"container %s of deployment %s does not have correct image: expected: %v got: %v",
						container.Name, obj.GetName(), expectedContainerImage, container.Image)

				}),
				pandariaAlertingDriversDeployExistsCheck,
			},
		},
		{

			Name: "Check affinity",
			Covers: []string{
				".Values.affinity",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					if obj.GetName() != "pandaria-alerting-drivers" {
						return
					}

					checker.MapSet(tc, PandariaAlertingDriversDeployExistsCheck, FoundKey, true)

					expectedAffinity, _ := checker.RenderValue[*corev1.Affinity](tc, ".Values.affinity")
					if expectedAffinity != nil && (*expectedAffinity) == (corev1.Affinity{}) {
						expectedAffinity = nil
					}

					assert.Equal(tc.T,
						expectedAffinity, podTemplateSpec.Spec.Affinity,
						"deployment %s does not have correct affinity: expected: %v, got: %v",
						obj.GetName(), expectedAffinity, podTemplateSpec.Spec.Affinity)

				}),
				pandariaAlertingDriversDeployExistsCheck,
			},
		},
		{

			Name: "Check resources",
			Covers: []string{
				".Values.resources",
			},
			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {

					if obj.GetName() != "pandaria-alerting-drivers" {
						return
					}

					checker.MapSet(tc, PandariaAlertingDriversDeployExistsCheck, FoundKey, true)

					ok := assert.Equal(tc.T, 1, len(podTemplateSpec.Spec.Containers),
						"deployment %s does not have correct number of container: expected: %d, got: %d",
						obj.GetName(), 1, len(podTemplateSpec.Spec.Containers))

					if !ok {
						return
					}

					container := podTemplateSpec.Spec.Containers[0]

					expectedResourceReq, _ := checker.RenderValue[corev1.ResourceRequirements](tc, ".Values.resources")

					assert.Equal(tc.T,
						expectedResourceReq, container.Resources,
						"container %s of deployment %s does not have correct resources constraint: expected: %v, got: %v",
						container.Name, obj.GetName(), expectedResourceReq, container.Resources)

				}),
				pandariaAlertingDriversDeployExistsCheck,
			},
		},
		{
			Name: "Check port",
			Covers: []string{
				".Values.port",
			},

			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, deployment *appsv1.Deployment) {

					expectedPort, exists := checker.RenderValue[int32](tc, ".Values.port")

					if exists {
						for _, container := range deployment.Spec.Template.Spec.Containers {

							assert.Equal(tc.T,
								expectedPort, container.Ports[0].ContainerPort,
								"container %s of deployment %s does not have correct port: expected: %v got: %v",
								container.Name, deployment.Name, expectedPort, container.Ports[0].ContainerPort)
						}
					}
				}),
			},
		},
		{
			Name: "Check replicas",
			Covers: []string{
				".Values.replicas",
			},

			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, deployment *appsv1.Deployment) {

					expectedReplicas, exists := checker.RenderValue[int32](tc, ".Values.replicas")

					if exists {
						assert.Equal(tc.T,
							expectedReplicas, *deployment.Spec.Replicas,
							"deployment %s does not have correct replicas: expected: %v got: %v",
							deployment.Name, expectedReplicas, *deployment.Spec.Replicas)
					}
				}),
			},
		},
	},
}

var pandariaAlertingDriversDeployExistsCheck = checker.Once(func(tc *checker.TestContext) {

	foundPandariaAlertingDriversDeploy, _ := checker.MapGet[string, string, bool](tc, PandariaAlertingDriversDeployExistsCheck, FoundKey)
	if !foundPandariaAlertingDriversDeploy {
		tc.T.Error("err: pandaria-alerting-drivers depoloyment not found")
	}
})
