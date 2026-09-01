{{- define "kopilot.subAgent.kubernetes.enabled" -}}
{{- if .Values.subAgents.kubernetes.enabled -}}true{{- end -}}
{{- end -}}

{{- define "kopilot.subAgent.helm.enabled" -}}
{{- if .Values.subAgents.helm.enabled -}}true{{- end -}}
{{- end -}}

{{- define "kopilot.subAgent.tavily.enabled" -}}
{{- if .Values.subAgents.tavily.enabled -}}true{{- end -}}
{{- end -}}

{{- define "kopilot.subAgent.rancher.enabled" -}}
{{- if .Values.subAgents.rancher.enabled -}}true{{- end -}}
{{- end -}}

{{- define "kopilot.subAgent.fleet.enabled" -}}
{{- if .Values.subAgents.fleet.enabled -}}true{{- end -}}
{{- end -}}

{{- define "kopilot.subAgent.provisioning.enabled" -}}
{{- if .Values.subAgents.provisioning.enabled -}}true{{- end -}}
{{- end -}}

{{- define "kopilot.deploy.kopilotMCPServer" -}}
{{- if or (include "kopilot.subAgent.kubernetes.enabled" .) (include "kopilot.subAgent.helm.enabled" .) (include "kopilot.subAgent.tavily.enabled" .) -}}true{{- end -}}
{{- end -}}

{{- define "kopilot.deploy.rancherMCPServer" -}}
{{- if or (include "kopilot.subAgent.rancher.enabled" .) (include "kopilot.subAgent.fleet.enabled" .) (include "kopilot.subAgent.provisioning.enabled" .) -}}true{{- end -}}
{{- end -}}

{{- define "kopilot.registry.default" -}}
{{- if .Values.global.cattle.systemDefaultRegistry -}}
{{- printf "%s/" .Values.global.cattle.systemDefaultRegistry -}}
{{- else -}}
{{- "registry.rancher.cn/" -}}
{{- end -}}
{{- end -}}

{{- define "kopilot.bundle.image" -}}
{{- printf "%s%s:%s" (include "kopilot.registry.default" .) .Values.bundle.image.repository .Values.bundle.image.tag -}}
{{- end -}}

{{- define "kopilot.runtime.image" -}}
{{- printf "%s%s:%s" (include "kopilot.registry.default" .) .Values.runtime.image.repository .Values.runtime.image.tag -}}
{{- end -}}

{{- define "kopilot.redis.name" -}}
{{- .Values.redis.name -}}
{{- end -}}

{{- define "kopilot.redis.pvc.name" -}}
{{- .Values.redis.persistentVolume.name -}}
{{- end -}}

{{- define "kopilot.redis.url" -}}
{{- if .Values.redis.externalURL -}}
    {{- .Values.redis.externalURL -}}
{{- else -}}
    {{- printf "%s.%s.svc.cluster.local:6379" (include "kopilot.redis.name" .) .Release.Namespace -}}
{{- end -}}
{{- end -}}

{{- define "kopilot.postgres.name" -}}
{{- .Values.postgres.name -}}
{{- end -}}

{{- define "kopilot.postgres.pvc.name" -}}
{{- .Values.postgres.persistentVolume.name -}}
{{- end -}}

{{- define "kopilot.postgres.password" -}}
{{- $password := .Values.postgres.password -}}
{{- $existingSecret := lookup "v1" "Secret" .Release.Namespace (include "kopilot.postgres.name" .) -}}
{{- if not $existingSecret -}}
  {{- $existingSecret = lookup "v1" "Secret" .Release.Namespace "postgres" -}}
{{- end -}}
{{- if $existingSecret -}}
  {{- if $existingSecret.data -}}
    {{- if $existingSecret.data.POSTGRES_PASSWORD -}}
      {{- $password = $existingSecret.data.POSTGRES_PASSWORD | b64dec -}}
    {{- end -}}
  {{- end -}}
{{- end -}}
{{- $password -}}
{{- end -}}

{{- define "kopilot.postgres.env" -}}
- name: POSTGRES_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "kopilot.postgres.name" . }}
      key: POSTGRES_PASSWORD
{{- end -}}

{{- define "kopilot.postgres.url" -}}
{{- if .Values.postgres.externalURL -}}
    {{- (splitList ":" .Values.postgres.externalURL) | first -}}
{{- else -}}
    {{- printf "%s.%s.svc.cluster.local" (include "kopilot.postgres.name" .) .Release.Namespace -}}
{{- end -}}
{{- end -}}

{{- define "kopilot.postgres.port" -}}
{{- if .Values.postgres.externalURL -}}
    {{- $spec := (splitList ":" .Values.postgres.externalURL) -}}
    {{- if gt (len $spec) 1 }}
        {{- index $spec 1 -}}
    {{- else -}}
        {{- 5432 -}}
    {{- end -}}
{{- else -}}
    {{- 5432 -}}
{{- end -}}
{{- end -}}

{{- define "kopilot.postgres.langgraph.url" -}}
{{- if .Values.postgres.externalURL -}}
    {{- printf "postgres://%s:%s@%s/langgraph" .Values.postgres.username .Values.postgres.password .Values.postgres.externalURL | b64enc | quote -}}
{{- else -}}
    {{- printf "postgres://%s:%s@%s.%s.svc.cluster.local:5432/langgraph" .Values.postgres.username .Values.postgres.password (include "kopilot.postgres.name" .) .Release.Namespace | b64enc | quote -}}
{{- end -}}
{{- end -}}

{{- define "kopilot.tavily.env.from" -}}
{{- if .Values.tavily.enabled }}
envFrom:
  - secretRef:
      name: kopilot-secret
{{- end -}}
{{- end -}}

{{- define "kopilot.tavily.api.key" -}}
{{- if .Values.tavily.enabled -}}
    {{- printf "%s" .Values.tavily.apiKey | trim | b64enc -}}
{{- end -}}
{{- end -}}

{{- define "kopilot.bootstrap.env" -}}
{{- if .Values.bootstrapPassword }}
- name: KOPILOT_BOOTSTRAP_TOKEN
  valueFrom:
    secretKeyRef:
      name: bootstrap-secret
      key: bootstrapPassword
{{- end -}}
{{- end -}}

{{- define "kopilot.rancher.sso.env" -}}
{{- if .Values.rancher_sso.enabled }}
- name: KOPILOT_RANCHER_SSO
  value: {{ .Values.rancher_sso.enabled | quote }}
- name: KOPILOT_RANCHER_SERVER_URL
  value: {{ .Values.rancher_sso.url | quote }}
{{- end -}}
{{- end -}}

{{- define "kopilot.privateCA.volume.mounts" -}}
{{- if .Values.privateCA }}
- mountPath: /etc/kopilot/ssl/cert.pem
  name: tls-cert-key-volume
  subPath: tls.crt
  readOnly: true
- mountPath: /etc/kopilot/ssl/key.pem
  name: tls-cert-key-volume
  subPath: tls.key
  readOnly: true
- mountPath: /etc/kopilot/ssl/cacerts.pem
  name: tls-ca-volume
  subPath: cacerts.pem
  readOnly: true
{{- end -}}
{{- end -}}

{{- define "kopilot.privateCA.volumes" -}}
{{- if .Values.privateCA }}
- name: tls-ca-volume
  secret:
    defaultMode: 0400
    secretName: tls-kopilot-ca
- name: tls-cert-key-volume
  secret:
    defaultMode: 0400
    secretName: tls-kopilot-ingress
{{- end -}}
{{- end -}}
