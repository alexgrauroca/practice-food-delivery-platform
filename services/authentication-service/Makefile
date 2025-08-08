openapi-validate-sync:
	@cd ../.. && ./scripts/openapi-validate-sync.sh services/authentication-service

openapi-validate-docs:
	@cd ../.. && ./scripts/openapi-validate-docs.sh services/authentication-service

openapi-generate:
	@cd ../.. && ./scripts/openapi-bundle.sh services/authentication-service

generate-clients: install-oapi-codegen
	@echo "Generating clients for authentication service..."
	@cd ../.. && mkdir -p clients/authentication-service
	@oapi-codegen -config client-gen-config.yaml docs/dist/openapi.yaml > ../../clients/authentication-service/client.go
	@cd ../../clients && go mod tidy

install-oapi-codegen:
	@bash ../../scripts/install-oapi-codegen.sh

generate-mocks: install-mockgen
	@echo "Generating mocks for authentication service..."
	@go generate ./...

install-mockgen:
	@bash ../../scripts/install-mockgen.sh
	@go get go.uber.org/mock/mockgen

run-unit-tests: generate-mocks
	@echo "Running unit tests for authentication service..."
	@go test -v -tags=unit ./...

run-integration-tests: start-mongodb
	@echo "Running integration tests for authentication service..."
	@bash -c "set -a && source ../../.env && set +a && go test -v -tags=integration ./..."

start-mongodb:
	@echo "Starting MongoDB Docker containers for integration tests..."
	@docker compose -f ../../docker-compose.yml up -d mongodb
