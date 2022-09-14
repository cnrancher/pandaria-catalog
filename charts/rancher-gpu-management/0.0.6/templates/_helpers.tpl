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

{{/*
Use new scheduler config
*/}}
{{- define "use_new_scheduler_config" -}}
{{- if eq .Values.image.defaultScheduler.version "v1.19" -}}
{{- "true" -}}
{{- else if .Values.image.defaultScheduler.version "v1.20" -}}
{{- "true" -}}
{{- else if .Values.image.defaultScheduler.version "v1.21" -}}
{{- "true" -}}
{{- else -}}
{{- "" -}}
{{- end -}}
{{- end -}}

