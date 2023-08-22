package rancher_k8s_auditlog_collector

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

var ChartPath = utils.MustGetPathFromModuleRoot("/charts/rancher-k8s-auditlog-collector/0.0.2")

var (
	DefaultReleaseName = "rancher-k8s-auditlog-collector-test-1"
	DefaultNamespace   = "default"
)

var suite = test.Suite{
	ChartPath: ChartPath,

	Cases: []test.Case{
		{
			Name: "File mode",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("auditlogServerPort", "9000").
				Set("auditlogServerHost", "1.1.1.1").
				Set("clusterid", "c-xxxx"),
		},
		{
			Name: "Webhook mode",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).
				Set("auditlogServerPort", "9000").
				Set("auditlogServerHost", "1.1.1.1").
				Set("mode", "webhook").
				Set("webhook.servicePort", "8787").
				Set("webhook.serviceNodePort", "30121").
				Set("clusterid", "c-xxxx"),
		},
	},

	NamedChecks: []test.NamedCheck{
		{
			Name:   "Check auditlog server and cluster config",
			Covers: []string{"auditlogServerHost", "auditlogServerPort", "clusterid"},
			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, configmap *corev1.ConfigMap) {
					auditlogServerPort, _ := checker.RenderValue[string](tc, ".Values.auditlogServerPort")
					auditlogServerHost, _ := checker.RenderValue[string](tc, ".Values.auditlogServerHost")
					clusterid, _ := checker.RenderValue[string](tc, ".Values.clusterid")

					assert.Contains(tc.T,
						configmap.Data, "main.conf",
						"%T %s does not have 'main.conf' key", configmap, checker.Key(configmap),
					)
					assert.Contains(tc.T,
						configmap.Data["main.conf"], auditlogServerHost,
						"%T %s does not have auditlog server host setting", configmap, checker.Key(configmap),
					)
					assert.Contains(tc.T,
						configmap.Data["main.conf"], auditlogServerPort,
						"%T %s does not have auditlog server port setting", configmap, checker.Key(configmap),
					)
					assert.Contains(tc.T,
						configmap.Data["main.conf"], clusterid,
						"%T %s does not have cluster id setting", configmap, checker.Key(configmap),
					)
				}),
			},
		},
		{
			Name: "Check fluentbit repository and tag",
			Covers: []string{
				".Values.global.systemDefaultRegistry",
				".Values.fluentbit.image.repository",
				".Values.fluentbit.image.tag",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetName() != "rancher-k8s-auditlog-fluentbit-file" &&
						obj.GetName() != "rancher-k8s-auditlog-fluentbit-webhook" {
						return
					}
					systemDefaultRegistry, _ := checker.RenderValue[string](tc, ".Values.global.systemDefaultRegistry")
					fluentbitRepository, _ := checker.RenderValue[string](tc, ".Values.fluentbit.image.repository")
					fluentbitTag, _ := checker.RenderValue[string](tc, ".Values.fluentbit.image.tag")
					expectedFluentbitImage := fluentbitRepository + ":" + fluentbitTag

					if systemDefaultRegistry != "" {
						expectedFluentbitImage = systemDefaultRegistry + "/" + expectedFluentbitImage
					}

					for _, container := range podTemplateSpec.Spec.Containers {
						assert.Equal(tc.T, expectedFluentbitImage, container.Image,
							"workload %s (type: %T) in Deployment %s/%s does not have correct image",
							container.Name, obj, obj.GetNamespace(), obj.GetName(),
						)
					}
				}),
			},
		},

		{
			Name: "Check webhook service",
			Covers: []string{
				".Values.webhook.servicePort",
				".Values.webhook.serviceNodePort",
			},

			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, service *corev1.Service) {
					mode, _ := checker.RenderValue[string](tc, ".Values.mode")
					if mode != "webhook" {
						return
					}

					servicePort, _ := checker.RenderValue[string](tc, ".Values.webhook.servicePort")
					serviceNodePort, _ := checker.RenderValue[string](tc, ".Values.webhook.serviceNodePort")

					for _, port := range service.Spec.Ports {
						assert.Equal(tc.T,
							fmt.Sprint(port.TargetPort.IntVal), servicePort,
							"%T %s does not serivceport error", service, checker.Key(service),
						)
						assert.Equal(tc.T,
							fmt.Sprint(port.NodePort), serviceNodePort,
							"%T %s does not serivceport error", service, checker.Key(service),
						)
					}

				}),
			},
		},

		{
			Name: "Check file mode",
			Covers: []string{
				".Values.apiserverLogPath",
				".Values.apiserverLogFile",
			},

			Checks: test.Checks{
				checker.PerResource(func(tc *checker.TestContext, configmap *corev1.ConfigMap) {
					mode, _ := checker.RenderValue[string](tc, ".Values.mode")
					if mode != "file" {
						return
					}

					apiserverLogFile, _ := checker.RenderValue[string](tc, ".Values.apiserverLogFile")

					assert.Contains(tc.T,
						configmap.Data["main.conf"], fmt.Sprintf("/var/log/kubernetes/audit/%s", apiserverLogFile),
						"%T %s does not serivceport error", configmap, checker.Key(configmap),
					)
				}),
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					mode, _ := checker.RenderValue[string](tc, ".Values.mode")
					if mode != "file" {
						return
					}

					if obj.GetName() != "rancher-k8s-auditlog-fluentbit-file" {
						return
					}
					apiserverLogPath, _ := checker.RenderValue[string](tc, ".Values.apiserverLogPath")

					for _, v := range podTemplateSpec.Spec.Volumes {
						if v.Name == "kubeapiserver-log" {
							assert.Equal(tc.T, apiserverLogPath, v.HostPath.Path,
								"%s/%s does not have correct volume",
								obj.GetNamespace(), obj.GetName(),
							)
						}
					}
				}),
			},
		},
	},
}
