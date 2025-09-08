{{/*
Expand the name of the chart.
*/}}
{{- define "tutuplapak.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "tutuplapak.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "tutuplapak.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "tutuplapak.labels" -}}
helm.sh/chart: {{ include "tutuplapak.chart" . }}
{{ include "tutuplapak.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "tutuplapak.selectorLabels" -}}
app.kubernetes.io/name: {{ include "tutuplapak.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "tutuplapak.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "tutuplapak.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
PostgreSQL fullname
*/}}
{{- define "tutuplapak.postgresql.fullname" -}}
{{- if .Values.database.postgresql.enabled }}
{{- printf "%s-postgresql" (include "tutuplapak.fullname" .) }}
{{- else }}
{{- .Values.database.postgresql.external.host }}
{{- end }}
{{- end }}

{{/*
MinIO fullname
*/}}
{{- define "tutuplapak.minio.fullname" -}}
{{- if .Values.minio.enabled }}
{{- printf "%s-minio" (include "tutuplapak.fullname" .) }}
{{- else }}
{{- .Values.minio.external.host }}
{{- end }}
{{- end }}

{{/*
Prometheus fullname
*/}}
{{- define "tutuplapak.prometheus.fullname" -}}
{{- if .Values.prometheus.enabled }}
{{- printf "%s-prometheus" (include "tutuplapak.fullname" .) }}
{{- else }}
{{- .Values.prometheus.external.host }}
{{- end }}
{{- end }}

{{/*
Grafana fullname
*/}}
{{- define "tutuplapak.grafana.fullname" -}}
{{- if .Values.grafana.enabled }}
{{- printf "%s-grafana" (include "tutuplapak.fullname" .) }}
{{- else }}
{{- .Values.grafana.external.host }}
{{- end }}
{{- end }}
