package rancher_gpu_sharing

import (
	"github.com/rancher/hull/pkg/chart"
	"github.com/rancher/hull/pkg/checker"
	"github.com/rancher/hull/pkg/test"
	"github.com/rancher/hull/pkg/utils"
	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ChartPath = utils.MustGetPathFromModuleRoot("/charts/rancher-gpu-sharing/0.0.1")

var (
	DefaultReleaseName = "rancher-sharing-test-1"
	DefaultNamespace   = "default"
)

var suite = test.Suite{
	ChartPath: ChartPath,

	Cases: []test.Case{
		{
			Name: "1.23",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("defaultScheduler.version", "v1.23"),
		},
		{
			Name: "1.24",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("defaultScheduler.version", "v1.24"),
		},
		{
			Name: "1.25",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("defaultScheduler.version", "v1.25"),
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
			Name: "Check shared device plugin repository and tag",
			Covers: []string{
				".Values.sharedeviceplugin.image.repository",
				".Values.sharedeviceplugin.image.tag",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "gpushare-device-plugin" {
						return
					}
					systemDefaultRegistry, _ := checker.RenderValue[string](tc, ".Values.global.systemDefaultRegistry")
					sdRepository, _ := checker.RenderValue[string](tc, ".Values.sharedeviceplugin.image.repository")
					sdTag, _ := checker.RenderValue[string](tc, ".Values.sharedeviceplugin.image.tag")
					expectedSDImage := sdRepository + ":" + sdTag

					if systemDefaultRegistry != "" {
						expectedSDImage = systemDefaultRegistry + "/" + expectedSDImage
					}

					for _, container := range podTemplateSpec.Spec.Containers {
						assert.Equal(tc.T, expectedSDImage, container.Image,
							"workload %s (type: %T) in Deployment %s/%s does not have correct image",
							container.Name, obj, obj.GetNamespace(), obj.GetName(),
						)
					}
				}),
			},
		},

		{
			Name: "Check defaultScheduler repository and tag",
			Covers: []string{
				".Values.defaultScheduler",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "gpushare-schd-extender" {
						return
					}
					systemDefaultRegistry, _ := checker.RenderValue[string](tc, ".Values.global.systemDefaultRegistry")

					defaultSchedulerVersion, _ := checker.RenderValue[string](tc, ".Values.defaultScheduler.version")

					supportedVersions, _ := checker.RenderValue[map[string]interface{}](tc,
						".Values.defaultScheduler.supportedVersions")

					image123 := supportedVersions["v1.23"].(map[string]interface{})
					defaultSchedulerRepository123 := image123["repository"].(string)
					defaultSchedulerTag123 := image123["tag"].(string)
					expectedImage123 := defaultSchedulerRepository123 + ":" + defaultSchedulerTag123

					image124 := supportedVersions["v1.24"].(map[string]interface{})
					defaultSchedulerRepository124 := image124["repository"].(string)
					defaultSchedulerTag124 := image124["tag"].(string)
					expectedImage124 := defaultSchedulerRepository124 + ":" + defaultSchedulerTag124

					schedulerExtenderRepo, _ := checker.RenderValue[string](tc, ".Values.schedulerextender.image.repository")
					schedulerExtenderTag, _ := checker.RenderValue[string](tc, ".Values.schedulerextender.image.tag")
					expectedSchedulerExtenderImage := schedulerExtenderRepo + ":" + schedulerExtenderTag

					if systemDefaultRegistry != "" {
						expectedImage123 = systemDefaultRegistry + "/" + expectedImage123
						expectedImage124 = systemDefaultRegistry + "/" + expectedImage124
						expectedSchedulerExtenderImage = systemDefaultRegistry + "/" + expectedSchedulerExtenderImage
					}

					for _, container := range podTemplateSpec.Spec.Containers {
						if container.Name == "gpushare-default-scheduler" {
							if defaultSchedulerVersion == "v1.23" {
								assert.Equal(tc.T, expectedImage123, container.Image,
									"workload %s (type: %T) in Deployment %s/%s does not have correct image",
									container.Name, obj, obj.GetNamespace(), obj.GetName(),
								)
							} else if defaultSchedulerVersion == "v1.24" {
								assert.Equal(tc.T, expectedImage124, container.Image,
									"workload %s (type: %T) in Deployment %s/%s does not have correct image",
									container.Name, obj, obj.GetNamespace(), obj.GetName(),
								)
							}
						} else if container.Name == "gpushare-scheduler-extender" {
							assert.Equal(tc.T, expectedSchedulerExtenderImage, container.Image,
								"workload %s (type: %T) in Deployment %s/%s does not have correct image",
								container.Name, obj, obj.GetNamespace(), obj.GetName(),
							)
						}
					}
				}),
			},
		},

		{
			Name: "Check runtime class",
			Covers: []string{
				".Values.defaultScheduler",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "gpushare-device-plugin" {
						return
					}

					runtimeClassName, _ := checker.RenderValue[string](tc, ".Values.runtimeClassName")
					if podTemplateSpec.Spec.RuntimeClassName == nil {
						assert.Equal(tc.T, runtimeClassName, "",
							"workload %s (type: %T) in Deployment %s/%s does not have correct image",
							podTemplateSpec.Name, obj, obj.GetNamespace(), obj.GetName(),
						)
					} else {
						assert.Equal(tc.T, runtimeClassName, *podTemplateSpec.Spec.RuntimeClassName,
							"workload %s (type: %T) in Deployment %s/%s does not have correct image",
							podTemplateSpec.Name, obj, obj.GetNamespace(), obj.GetName(),
						)
					}

				}),
			},
		},
	},
}
