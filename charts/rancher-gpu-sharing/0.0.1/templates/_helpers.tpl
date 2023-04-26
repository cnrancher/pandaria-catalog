{{/*
Get default scheduler image version by kubernetes version
Supported versions map is set in .Values.image.defaultScheduler.supportedVersions
*/}}
{{- define "gpushare.defaultscheduler.image" -}}
{{- range $key, $val := .Values.defaultScheduler.supportedVersions }}
{{- if eq $.Values.defaultScheduler.version $key -}}
    {{- printf "%s:%s" $val.repository $val.tag -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Get Rancher system-default-registry
*/}}
{{- define "system_default_registry" -}}
{{- if .Values.global.systemDefaultRegistry -}}
{{- printf "%s/" .Values.global.systemDefaultRegistry -}}
{{- else -}}
{{- "" -}}
{{- end -}}
{{- end -}}

{{/*
Scheduler config apiversion
*/}}
{{- define "scheduler_extender_apiversion" -}}
{{- $vbeta1list := list "v1.19" "v1.20" "v1.21" -}}
{{- if has .Values.defaultScheduler.version $vbeta1list -}}
{{- "kubescheduler.config.k8s.io/v1beta1" -}}
{{- else if eq .Values.defaultScheduler.version "v1.22" -}}
{{- "kubescheduler.config.k8s.io/v1beta2" -}}
{{- else -}}
{{- "kubescheduler.config.k8s.io/v1beta3" -}}
{{- end -}}
{{- end -}}
