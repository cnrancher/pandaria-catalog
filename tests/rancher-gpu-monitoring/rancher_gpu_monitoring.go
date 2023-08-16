package rancher_gpu_monitoring

import (
	"github.com/rancher/hull/pkg/chart"
	"github.com/rancher/hull/pkg/checker"
	"github.com/rancher/hull/pkg/test"
	"github.com/rancher/hull/pkg/utils"
	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ChartPath = utils.MustGetPathFromModuleRoot("/charts/rancher-gpu-monitoring/0.0.1")

var (
	DefaultReleaseName = "rancher-monitoring-test-1"
	DefaultNamespace   = "default"
)

var suite = test.Suite{
	ChartPath: ChartPath,

	Cases: []test.Case{
		{
			Name: "Default",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace),
		},
		{
			Name: "systemdefaultregistry",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("global.systemDefaultRegistry", "testrancher.io"),
		},
		{
			Name: "Use runtime class",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("runtimeClassName", "nvidia"),
		},
	},

	NamedChecks: []test.NamedCheck{
		{
			Name: "Check image repository and tag",
			Covers: []string{
				".Values.image.repository",
				".Values.image.tag",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					systemDefaultRegistry, _ := checker.RenderValue[string](tc, ".Values.global.systemDefaultRegistry")
					imageRepository, _ := checker.RenderValue[string](tc, ".Values.image.repository")
					imageTag, _ := checker.RenderValue[string](tc, ".Values.image.tag")
					expectedImage := imageRepository + ":" + imageTag

					if systemDefaultRegistry != "" {
						expectedImage = systemDefaultRegistry + "/" + expectedImage
					}

					for _, container := range podTemplateSpec.Spec.Containers {
						assert.Equal(tc.T, expectedImage, container.Image,
							"workload %s (type: %T) in workload %s/%s does not have correct image",
							container.Name, obj, obj.GetNamespace(), obj.GetName(),
						)
					}
				}),
			},
		},

		{
			Name: "Check runtime class",
			Covers: []string{
				".Values.runtimeClassName",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					runtimeClassName, _ := checker.RenderValue[string](tc, ".Values.runtimeClassName")
					if podTemplateSpec.Spec.RuntimeClassName == nil {
						assert.Equal(tc.T, runtimeClassName, "",
							"workload %s (type: %T) in Deployment %s/%s does not have correct runtimeclass",
							podTemplateSpec.Name, obj, obj.GetNamespace(), obj.GetName(),
						)
					} else {
						assert.Equal(tc.T, runtimeClassName, *podTemplateSpec.Spec.RuntimeClassName,
							"workload %s (type: %T) in Deployment %s/%s does not have correct runtimeclass",
							podTemplateSpec.Name, obj, obj.GetNamespace(), obj.GetName(),
						)
					}

				}),
			},
		},
	},
}
