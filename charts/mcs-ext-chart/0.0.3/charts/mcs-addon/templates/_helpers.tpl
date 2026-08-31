{{- define "mcs-addon.component_name" -}}
{{- printf "%s-%s" .Chart.Name .componentName | trunc 63 | trimSuffix "-" | quote -}}
{{- end -}}

{{- define "mcs-addon.component_type_name" -}}
{{- printf "%s-%s-%s" .Chart.Name .componentName .Values.global.installationType | trunc 63 | trimSuffix "-" | quote -}}
{{- end -}}

{{- define "mcs-addon.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" | quote -}}
{{- end -}}

{{- define "mcs-addon.chart_labels" -}}
heritage: {{ .Release.Service | quote }}
release: {{ .Release.Name | quote }}
chart: {{ template "mcs-addon.chart" . }}
app: {{ .Chart.Name | trunc 63 | trimSuffix "-" | quote -}}
{{- end -}}

{{- define "mcs-addon.broker_role_name" -}}
{{- printf "%s-%s" .Chart.Name "broker" | trunc 63 | trimSuffix "-" | quote -}}
{{- end -}}

{{- define "system_default_registry" -}}
{{- if .Values.global.cattle.systemDefaultRegistry -}}
{{- printf "%s/" .Values.global.cattle.systemDefaultRegistry -}}
{{- end -}}
{{- end -}}
