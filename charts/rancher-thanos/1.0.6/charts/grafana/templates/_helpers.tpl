{{/* vim: set filetype=mustache: */}}

{{- define "app.dashboards.fullname" -}}
{{- printf "%s-dashboards" .Release.Name -}}
{{- end -}}

{{- define "app.dashboards-provisionings.fullname" -}}
{{- printf "%s-dashboards-provisionings" .Release.Name -}}
{{- end -}}

{{- define "app.provisionings.fullname" -}}
{{- printf "%s-provisionings" .Release.Name -}}
{{- end -}}

{{- define "app.dashboards.istio-fullname" -}}
{{- $name := include "app.name" . -}}
{{- printf "%s-%s-dashboards" "istio" .Release.Name -}}
{{- end -}}


{{- define "system_default_registry" -}}
{{- if .Values.global.systemDefaultRegistry -}}
{{- printf "%s/" .Values.global.systemDefaultRegistry -}}
{{- else -}}
{{- "" -}}
{{- end -}}
{{- end -}}
