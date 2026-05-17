MAIN        := cmd/app/main.go
DOCS_DIR    := docs

.PHONY: help swagger swagger-fmt run build

help:
	@echo "Targets:"
	@echo "  swagger      Generate OpenAPI docs ($(DOCS_DIR)/) from handler annotations"
	@echo "  swagger-fmt  Format swag comments in Go sources"
	@echo "  run          Run the service"
	@echo "  build        Build binary to bin/course-service"

swagger:
	swag init -g ./internal/transport/http/server.go -o ./docs --parseDependency --parseInternal

swagger-fmt:
	swag fmt -g ./internal/transport/http/server.go

run:
	go run ./cmd/app

build:
	go build -o bin/course-service ./cmd/app
