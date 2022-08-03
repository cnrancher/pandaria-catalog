{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}

{{- define "system_default_registry" -}}
{{- if .Values.systemDefaultRegistry -}}
{{- printf "%s/" .Values.systemDefaultRegistry -}}
{{- else -}}
{{- "" -}}
{{- end -}}
{{- end -}}
