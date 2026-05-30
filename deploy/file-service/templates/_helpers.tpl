{{/*
=================================================================
_helpers.tpl — вспомогательные шаблоны (named templates).
Используются другими файлами через `include "course-service.xxx" .`
Это стандартный паттерн для Helm: имена ресурсов, метки и т.д.
формируются в одном месте — иначе при ребилде имени в одном
из ресурсов получим расхождения, и k8s "потеряет" связь.
=================================================================
*/}}

{{/*
Базовое имя приложения. По умолчанию — name из Chart.yaml,
но можно перебить через .Values.nameOverride.
*/}}
{{- define "file-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Полное имя релиза, с которым именуются все ресурсы.
Логика по официальному шаблону Helm:
  - если задан fullnameOverride — берём его
  - если Release.Name уже содержит Chart.Name — используем как есть
  - иначе склеиваем Release.Name-Chart.Name
DNS-имена в k8s не длиннее 63 символов, поэтому trunc 63.
*/}}
{{- define "file-service.fullname" -}}
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
Идентификатор Chart'а в формате name-version (используется в labels).
*/}}
{{- define "file-service.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Стандартные метки (labels), которые ставятся на ВСЕ ресурсы.
Помогают потом смотреть `kubectl get all -l app.kubernetes.io/instance=<release>`.
*/}}
{{- define "file-service.labels" -}}
helm.sh/chart: {{ include "file-service.chart" . }}
{{ include "file-service.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Селекторные метки. Это ПОДМНОЖЕСТВО labels, по которым Service ищет Pod'ы,
а Deployment связывает свои Pod'ы со своим ReplicaSet.
ВАЖНО: их нельзя менять между релизами — k8s не разрешит обновить selector.
Поэтому здесь только name + instance (стабильные значения).
*/}}
{{- define "file-service.selectorLabels" -}}
app.kubernetes.io/name: {{ include "file-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Имя ServiceAccount'а: либо берём явное из values, либо генерируем по fullname.
*/}}
{{- define "file-service.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "file-service.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
Имя Secret'а с .env файлом.
Если пользователь подсунул свой existingSecret — используем его,
иначе создадим свой <fullname>-env.
*/}}
{{- define "file-service.envSecretName" -}}
{{- if .Values.secrets.existingSecret -}}
{{- .Values.secrets.existingSecret -}}
{{- else -}}
{{- printf "%s-env" (include "file-service.fullname" .) -}}
{{- end -}}
{{- end -}}
