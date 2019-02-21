{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "whitelister.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" | lower -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "whitelister.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "whitelister.labels.selector" -}}
app: {{ template "whitelister.name" . }}
group: {{ .Values.whitelister.labels.group }}
provider: {{ .Values.whitelister.labels.provider }}
{{- end -}}

{{- define "whitelister.labels.stakater" -}}
{{ template "whitelister.labels.selector" . }}
version: {{ .Values.whitelister.labels.version }}
{{- end -}}

{{- define "whitelister.labels.chart" -}}
chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
release: {{ .Release.Name | quote }}
heritage: {{ .Release.Service | quote }}
{{- end -}}