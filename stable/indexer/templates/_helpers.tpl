{{/* vim: set filetype=mustache: */}}
{{/*
Allow users to upgrade container versions as needed.
*/}}
{{- define "indexer.version" -}}
{{- default .Chart.AppVersion .Values.image.tag -}}
{{- end }}

{{/*
Expand the name of the chart.
*/}}
{{- define "indexer.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "indexer.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "indexer.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "indexer.labels" -}}
helm.sh/chart: {{ include "indexer.chart" . }}
{{ include "indexer.selectorLabels" . }}
app.kubernetes.io/version: {{ include "indexer.version" . | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: deps-cloud
app.kubernetes.io/component: indexer
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "indexer.selectorLabels" -}}
app.kubernetes.io/name: {{ include "indexer.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "indexer.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "indexer.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}
