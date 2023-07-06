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
