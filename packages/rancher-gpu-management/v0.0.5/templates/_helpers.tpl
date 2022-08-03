{{/*
Get default scheduler image version by kubernetes version
Supported versions map is set in .Values.image.defaultScheduler.supportedVersions
*/}}
{{- define "gpushare.defaultscheduler.image" -}}
{{- range $key, $val := .Values.image.defaultScheduler.supportedVersions }}
{{- if eq $.Values.image.defaultScheduler.version $key -}}
    {{ $val }}
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
