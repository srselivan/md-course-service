FROM golang:1.26.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cd cmd/app/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/course-service

FROM alpine:3.23
# UID/GID 1000 — совпадают с runAsUser/fsGroup в Helm chart.
# Важно: USER должен быть ЧИСЛОВЫМ, иначе при runAsNonRoot: true Kubernetes
# не сможет проверить, что контейнер не root (ошибка "non-numeric user (app)").
RUN apk add --no-cache ca-certificates \
    && addgroup -S -g 1000 app \
    && adduser -S -u 1000 -G app app
WORKDIR /app
COPY --from=builder /app/course-service .
COPY migrations ./migrations
RUN mkdir -p /app/logs && chown -R 1000:1000 /app
USER 1000:1000
EXPOSE 10001
CMD ["./course-service"]