{{/*
generate dingtalk options
*/}}
{{- define "dingtalk.options" -}}
  {{ $options := "" }}
  {{- if .Values.sink.dingtalk.label -}}
    {{ $options = printf "%s&label=%s" $options (toString .Values.sink.dingtalk.label)}}
  {{- end -}}
  {{- if .Values.sink.dingtalk.level -}}
    {{ $options = printf "%s&level=%s" $options (toString .Values.sink.dingtalk.level)}}
  {{- end -}}
  {{- if .Values.sink.dingtalk.namespaces -}}
    {{ $options = printf "%s&namespaces=%s" $options (toString .Values.sink.dingtalk.namespaces)}}
  {{- end -}}
  {{- if .Values.sink.dingtalk.kinds -}}
    {{ $options = printf "%s&kinds=%s" $options (toString .Values.sink.dingtalk.kinds)}}
  {{- end -}}
  {{- if .Values.sink.dingtalk.msg_type -}}
    {{ $options = printf "%s&msg_type=%s" $options (toString .Values.sink.dingtalk.msg_type)}}
  {{- end -}}
  {{- if .Values.sink.dingtalk.cluster_id -}}
    {{ $options = printf "%s&cluster_id=%s" $options (toString .Values.sink.dingtalk.cluster_id)}}
  {{- end -}}
  {{- if .Values.sink.dingtalk.region -}}
    {{ $options = printf "%s&region=%s" $options (toString .Values.sink.dingtalk.region)}}
  {{- end -}}

  {{- $options = (trimPrefix "&" $options) -}}

  {{- if contains "?" .Values.sink.dingtalk.webhook_url -}}
    {{- $options = printf "%s&%s" .Values.sink.dingtalk.webhook_url $options }}
  {{- else if gt (len $options) 0 }}
    {{- $options = printf "%s?%s" .Values.sink.dingtalk.webhook_url $options }}
  {{- else -}}
    {{- $options = printf "%s" .Values.sink.dingtalk.webhook_url }}
  {{- end -}}

  {{- $options }}
{{- end -}}


{{/*
generate sls options
*/}}
{{- define "sls.options" -}}
  {{ $options := "" }}
  {{- if .Values.sink.sls.project -}}
    {{ $options = printf "%s&project=%s" $options (toString .Values.sink.sls.project)}}
  {{- end -}}
  {{- if .Values.sink.sls.logStore -}}
    {{ $options = printf "%s&logStore=%s" $options (toString .Values.sink.sls.logStore)}}
  {{- end -}}
  {{- if .Values.sink.sls.topic -}}
    {{ $options = printf "%s&topic=%s" $options (toString .Values.sink.sls.topic)}}
  {{- end -}}

  {{ $options = (trimPrefix "&" $options) }}

  {{- if contains "?" .Values.sink.sls.sls_endpoint -}}
    {{- $options = printf "%s&%s" .Values.sink.sls.sls_endpoint $options }}
  {{- else if gt (len $options) 0 }}
    {{- $options = printf "%s?%s" .Values.sink.sls.sls_endpoint $options }}
  {{- else -}}
    {{- $options = printf "%s" .Values.sink.sls.sls_endpoint }}
  {{- end -}}

  {{- $options }}
{{- end -}}

{{/*
generate elasticsearch options
*/}}
{{- define "elasticsearch.options" -}}
  {{ $options := "" }}
  {{- if .Values.sink.elasticsearch.index -}}
    {{ $options = printf "%s&index=%s" $options (toString .Values.sink.elasticsearch.index)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.esUserName -}}
    {{ $options = printf "%s&esUserName=%s" $options (toString .Values.sink.elasticsearch.esUserName)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.esUserSecret -}}
    {{ $options = printf "%s&esUserSecret=%s" $options (toString .Values.sink.elasticsearch.esUserSecret)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.maxRetries -}}
    {{ $options = printf "%s&maxRetries=%s" $options (toString .Values.sink.elasticsearch.maxRetries)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.healthCheck -}}
    {{ $options = printf "%s&healthCheck=%s" $options (toString .Values.sink.elasticsearch.healthCheck)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.sniff -}}
    {{ $options = printf "%s&sniff=%s" $options (toString .Values.sink.elasticsearch.sniff)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.startupHealthcheckTimeout -}}
    {{ $options = printf "%s&startupHealthcheckTimeout=%s" $options (toString .Values.sink.elasticsearch.startupHealthcheckTimeout)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.ver -}}
    {{ $options = printf "%s&ver=%s" $options (toString .Values.sink.elasticsearch.ver)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.bulkWorkers -}}
    {{ $options = printf "%s&bulkWorkers=%s" $options (toString .Values.sink.elasticsearch.bulkWorkers)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.cluster_name -}}
    {{ $options = printf "%s&cluster_name=%s" $options (toString .Values.sink.elasticsearch.cluster_name)}}
  {{- end -}}
  {{- if .Values.sink.elasticsearch.pipeline -}}
    {{ $options = printf "%s&pipeline=%s" $options (toString .Values.sink.elasticsearch.pipeline)}}
  {{- end -}}

  {{ $options = (trimPrefix "&" $options) }}

  {{- if contains "?" .Values.sink.elasticsearch.es_server_url -}}
    {{- $options = printf "%s&%s" .Values.sink.elasticsearch.es_server_url $options }}
  {{- else if gt (len $options) 0  }}
    {{- $options = printf "%s?%s" .Values.sink.elasticsearch.es_server_url $options }}
  {{- else -}}
    {{- $options = printf "%s" .Values.sink.elasticsearch.es_server_url }}
  {{- end -}}

  {{- $options }}
{{- end -}}

{{/*
generate honeycomb options
*/}}
{{- define "honeycomb.options" -}}
  {{ $options := "" }}
  {{- if .Values.sink.honeycomb.dataset -}}
    {{ $options = printf "%s&dataset=%s" $options (toString .Values.sink.honeycomb.dataset)}}
  {{- end -}}
  {{- if .Values.sink.honeycomb.writekey -}}
    {{ $options = printf "%s&writekey=%s" $options (toString .Values.sink.honeycomb.writekey)}}
  {{- end -}}
  {{- if .Values.sink.honeycomb.apihost -}}
    {{ $options = printf "%s&apihost=%s" $options (toString .Values.sink.honeycomb.apihost)}}
  {{- end -}}

  {{ $options = printf "?%s" (trimPrefix "&" $options) }}

  {{- $options }}
{{- end -}}

{{/*
generate influxdb options
*/}}
{{- define "influxdb.options" -}}
  {{ $options := "" }}
  {{- if .Values.sink.influxdb.user -}}
    {{ $options = printf "%s&user=%s" $options (toString .Values.sink.influxdb.user)}}
  {{- end -}}
  {{- if .Values.sink.influxdb.pw -}}
    {{ $options = printf "%s&pw=%s" $options (toString .Values.sink.influxdb.pw)}}
  {{- end -}}
  {{- if .Values.sink.influxdb.db -}}
    {{ $options = printf "%s&db=%s" $options (toString .Values.sink.influxdb.db)}}
  {{- end -}}
  {{- if .Values.sink.influxdb.insecuressl -}}
    {{ $options = printf "%s&insecuressl=%s" $options (toString .Values.sink.influxdb.insecuressl)}}
  {{- end -}}
  {{- if .Values.sink.influxdb.withfields -}}
    {{ $options = printf "%s&withfields=%s" $options (toString .Values.sink.influxdb.withfields)}}
  {{- end -}}
  {{- if .Values.sink.influxdb.cluster_name -}}
    {{ $options = printf "%s&cluster_name=%s" $options (toString .Values.sink.influxdb.cluster_name)}}
  {{- end -}}


  {{ $options = (trimPrefix "&" $options) }}

  {{- if contains "?" .Values.sink.influxdb.influxdb_url -}}
    {{- $options = printf "%s&%s" .Values.sink.influxdb.influxdb_url $options }}
  {{- else if gt (len $options) 0  }}
    {{- $options = printf "%s?%s" .Values.sink.influxdb.influxdb_url $options }}
  {{- else -}}
    {{- $options = printf "%s" .Values.sink.influxdb.influxdb_url }}
  {{- end -}}

  {{- $options }}
{{- end -}}


{{/*
generate kafka options
*/}}
{{- define "kafka.options" -}}
  {{ $options := "" }}
  {{- if .Values.sink.kafka.brokers -}}
    {{ $options = printf "%s&%s" $options (toString .Values.sink.kafka.brokers)}}
  {{- end -}}
  {{- if .Values.sink.kafka.eventstopic -}}
    {{ $options = printf "%s&eventstopic=%s" $options (toString .Values.sink.influxdb.eventstopic)}}
  {{- end -}}
  {{- if .Values.sink.kafka.compression -}}
    {{ $options = printf "%s&compression=%s" $options (toString .Values.sink.influxdb.compression)}}
  {{- end -}}
  {{- if .Values.sink.kafka.user -}}
    {{ $options = printf "%s&user=%s" $options (toString .Values.sink.influxdb.user)}}
  {{- end -}}
  {{- if .Values.sink.kafka.password -}}
    {{ $options = printf "%s&password=%s" $options (toString .Values.sink.influxdb.password)}}
  {{- end -}}
  {{- if .Values.sink.kafka.cacert -}}
    {{ $options = printf "%s&cacert=%s" $options (toString .Values.sink.influxdb.cacert)}}
  {{- end -}}
  {{- if .Values.sink.kafka.cert -}}
    {{ $options = printf "%s&cert=%s" $options (toString .Values.sink.influxdb.cert)}}
  {{- end -}}
  {{- if .Values.sink.kafka.key -}}
    {{ $options = printf "%s&key=%s" $options (toString .Values.sink.influxdb.key)}}
  {{- end -}}
  {{- if .Values.sink.kafka.insecuressl -}}
    {{ $options = printf "%s&insecuressl=%s" $options (toString .Values.sink.influxdb.insecuressl)}}
  {{- end -}}

  {{ $options = printf "?%s" (trimPrefix "&" $options) }}

  {{- $options }}
{{- end -}}

{{/*
generate wechat options
*/}}
{{- define "wechat.options" -}}
  {{ $options := "" }}
  {{- if .Values.sink.wechat.corp_id -}}
    {{ $options = printf "%s&corp_id=%s" $options (toString .Values.sink.wechat.corp_id)}}
  {{- end -}}
  {{- if .Values.sink.wechat.corp_secret -}}
    {{ $options = printf "%s&corp_secret=%s" $options (toString .Values.sink.wechat.corp_secret)}}
  {{- end -}}
  {{- if .Values.sink.wechat.agent_id -}}
    {{- if (typeIs "float64" .Values.sink.wechat.agent_id) -}}
      {{ $options = printf "%s&agent_id=%s" $options (printf "%.0f" .Values.sink.wechat.agent_id )}}
    {{- else -}}
      {{ $options = printf "%s&agent_id=%s" $options (toString .Values.sink.wechat.agent_id )}}
    {{- end -}}
  {{- end -}}
  {{- if .Values.sink.wechat.to_user -}}
    {{ $options = printf "%s&to_user=%s" $options (toString .Values.sink.wechat.to_user)}}
  {{- end -}}
  {{- if .Values.sink.wechat.label -}}
    {{ $options = printf "%s&label=%s" $options (toString .Values.sink.wechat.label)}}
  {{- end -}}
  {{- if .Values.sink.wechat.level -}}
    {{ $options = printf "%s&level=%s" $options (toString .Values.sink.wechat.level)}}
  {{- end -}}
  {{- if .Values.sink.wechat.namespaces -}}
    {{ $options = printf "%s&namespaces=%s" $options (toString .Values.sink.wechat.namespaces)}}
  {{- end -}}
  {{- if .Values.sink.wechat.kinds -}}
    {{ $options = printf "%s&kinds=%s" $options (toString .Values.sink.wechat.kinds)}}
  {{- end -}}

  {{ $options = printf "?%s" (trimPrefix "&" $options) }}

  {{- $options }}
{{- end -}}

{{/*
generate webhook options
*/}}
{{- define "webhook.options" -}}
  {{ $options := "" }}
  {{- if .Values.sink.webhook.level -}}
    {{ $options = printf "%s&level=%s" $options (toString .Values.sink.webhook.level)}}
  {{- end -}}
  {{- if .Values.sink.webhook.namespaces -}}
    {{ $options = printf "%s&namespaces=%s" $options (toString .Values.sink.webhook.namespaces)}}
  {{- end -}}
  {{- if .Values.sink.webhook.kinds -}}
    {{ $options = printf "%s&kinds=%s" $options (toString .Values.sink.webhook.kinds)}}
  {{- end -}}
  {{- if .Values.sink.webhook.reason -}}
    {{ $options = printf "%s&reason=%s" $options (toString .Values.sink.webhook.reason)}}
  {{- end -}}
  {{- if .Values.sink.webhook.method -}}
    {{ $options = printf "%s&method=%s" $options (toString .Values.sink.webhook.method)}}
  {{- end -}}
  {{- if .Values.sink.webhook.header -}}
    {{ $options = printf "%s&header=%s" $options (toString .Values.sink.webhook.header)}}
  {{- end -}}
  {{- if .Values.sink.webhook.use_custom_body -}}
    {{- $options = printf "%s&custom_body_configmap=%s" $options (toString .Values.sink.webhook.custom_body_configmap_name) -}}
    {{- $options = printf "%s&custom_body_configmap_namespace=%s" $options (toString .Release.Namespace) -}}
  {{- end -}}

  {{ $options = (trimPrefix "&" $options) }}

  {{- if contains "?" .Values.sink.webhook.webhook_url -}}
    {{- $options = printf "%s&%s" .Values.sink.webhook.webhook_url $options }}
  {{- else if gt (len $options) 0  }}
    {{- $options = printf "%s?%s" .Values.sink.webhook.webhook_url $options }}
  {{- else -}}
    {{- $options = printf "%s" .Values.sink.webhook.webhook_url }}
  {{- end -}}

  {{- $options }}
{{- end -}}

{{/*
generate mysql options
*/}}
{{- define "mysql.options" -}}
  {{ $options := "" }}
  {{- $options = printf "%s" .Values.sink.mysql.mysql_jdbc_url }}
  {{- $options }}
{{- end -}}

{{/*
generate sink configuration
*/}}
{{- define "sink.config" -}}
  {{- if eq .Values.sinktarget "mysql" -}}
    {{ printf "--sink=mysql:?%s" (include "mysql.options" .) }}
  {{- else if eq .Values.sinktarget "kafka" -}}
    {{ printf "--sink=kafka:%s" (include "kafka.options" .) }}
  {{- else if eq .Values.sinktarget "elasticsearch" -}}
    {{ printf "--sink=elasticsearch:%s" (include "elasticsearch.options" .) }}
  {{- else if eq .Values.sinktarget "dingtalk" -}}
    {{ printf "--sink=dingtalk:%s" (include "dingtalk.options" .) }}
  {{- else if eq .Values.sinktarget "wechat" -}}
    {{ printf "--sink=wechat:%s" (include "wechat.options" .) }}
  {{- else if eq .Values.sinktarget "webhook" -}}
    {{ printf "--sink=webhook:%s" (include "webhook.options" .) }}
  {{- else if eq .Values.sinktarget "influxdb" -}}
    {{ printf "--sink=influxdb:%s" (include "influxdb.options" .) }}
  {{- else if eq .Values.sinktarget "honeycomb" -}}
    {{ printf "--sink=honeycomb:%s" (include "honeycomb.options" .) }}
  {{- else if eq .Values.sinktarget "sls" -}}
    {{ printf "--sink=sls:%s" (include "sls.options" .) }}
  {{- end -}}
{{- end -}}
