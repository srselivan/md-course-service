FROM golang:1.26.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cd cmd/app/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/course-service

FROM alpine:3.23
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/course-service .
EXPOSE 10001
CMD ["./course-service"]