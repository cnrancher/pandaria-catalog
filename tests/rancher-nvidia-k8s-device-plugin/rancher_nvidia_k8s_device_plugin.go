package rancher_nvidia_k8s_device_plugin

import (
	"fmt"

	"github.com/rancher/hull/pkg/chart"
	"github.com/rancher/hull/pkg/checker"
	"github.com/rancher/hull/pkg/test"
	"github.com/rancher/hull/pkg/utils"
	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ChartPath = utils.MustGetPathFromModuleRoot("/charts/rancher-nvidia-k8s-device-plugin/0.0.3")

var (
	DefaultReleaseName = "rancher-nvidia-k8s-device-plugin-test-1"
	DefaultNamespace   = "test-nvidia"
)

var suite = test.Suite{
	ChartPath: ChartPath,

	Cases: []test.Case{
		{
			Name: "Default",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("global.systemDefaultRegistry", "testrancher.io"),
		},
		{
			Name: "systemdefaultregistry",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("global.systemDefaultRegistry", "testrancher.io"),
		},
		{
			Name: "gfd disable",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("gfd.enabled", false).
				Set("global.systemDefaultRegistry", "testrancher.io"),
		},
		{
			Name: "Use runtime class",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("runtimeClassName", "nvidia").
				Set("global.systemDefaultRegistry", "testrancher.io"),
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
						if container.Name == "nvidia-device-plugin-init" ||
							container.Name == "nvidia-device-plugin-sidecar" ||
							container.Name == "nvidia-device-plugin-ctr" {
							assert.Equal(tc.T, expectedImage, container.Image,
								"workload %s (type: %T) in workload %s/%s does not have correct image",
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
				".Values.runtimeClassName",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != DefaultReleaseName {
						return
					}

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

		{
			Name: "Check gfd image",
			Covers: []string{
				".Values.gfd.image.repository",
				".Values.gfd.image.tag",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					gfdenabled, _ := checker.RenderValue[bool](tc, ".Values.gfd.enabled")
					if !gfdenabled {
						return
					}

					systemDefaultRegistry, _ := checker.RenderValue[string](tc, ".Values.global.systemDefaultRegistry")
					imageRepository, _ := checker.RenderValue[string](tc, ".Values.gfd.image.repository")
					imageTag, _ := checker.RenderValue[string](tc, ".Values.gfd.image.tag")
					expectedImage := imageRepository + ":" + imageTag

					if systemDefaultRegistry != "" {
						expectedImage = systemDefaultRegistry + "/" + expectedImage
					}

					for _, container := range podTemplateSpec.Spec.Containers {
						if container.Name == "gpu-feature-discovery-init" ||
							container.Name == "gpu-feature-discovery-sidecar" ||
							container.Name == "gpu-feature-discovery-ctr" {
							assert.Equal(tc.T, expectedImage, container.Image,
								"workload %s (type: %T) in workload %s/%s does not have correct image",
								container.Name, obj, obj.GetNamespace(), obj.GetName(),
							)
						}
					}
				}),
			},
		},

		{
			Name: "Check gfd disable",
			Covers: []string{
				".Values.gfd.enbaled",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != DefaultReleaseName {
						return
					}
					gfdenabled, _ := checker.RenderValue[bool](tc, ".Values.gfd.enabled")
					if gfdenabled {
						return
					}
					fmt.Println(obj.GetName())
					hasNodeSelector := false
					for k, v := range podTemplateSpec.Spec.NodeSelector {
						if k == "gpu.cattle.io/type" && v == "nvidia" {
							hasNodeSelector = true
						}
					}

					assert.Equal(tc.T, true, hasNodeSelector,
						"workload %s/%s does not have nodeselector gpu.cattle.io/type:nvidia",
						obj.GetNamespace(), obj.GetName(),
					)
				}),
			},
		},
	},
}
