openapi-validate-sync:
	@cd ../.. && ./scripts/openapi-validate-sync.sh services/customer-service

openapi-validate-docs:
	@cd ../.. && ./scripts/openapi-validate-docs.sh services/customer-service

openapi-generate:
	@cd ../.. && ./scripts/openapi-bundle.sh services/customer-service

generate-clients: install-oapi-codegen
	@echo "Generating clients for customer service..."
	@cd ../.. && mkdir -p clients/customer-service
	@oapi-codegen -config client-gen-config.yaml docs/dist/openapi.yaml > ../../clients/customer-service/client.go
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
