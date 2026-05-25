FROM golang:1.26.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cd cmd/app/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/course-service

FROM alpine:3.23
RUN apk add --no-cache ca-certificates \
    && addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=builder /app/course-service .
COPY migrations ./migrations
RUN mkdir -p /app/logs && chown -R app:app /app
USER app
EXPOSE 10001
CMD ["./course-service"]