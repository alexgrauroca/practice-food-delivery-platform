openapi-validate-sync:
	@cd ../.. && ./scripts/openapi-validate-sync.sh services/customer-service

openapi-validate-docs:
	@cd ../.. && ./scripts/openapi-validate-docs.sh services/customer-service

openapi-generate:
	@cd ../.. && ./scripts/openapi-bundle.sh services/customer-service

generate-mocks: install-mockgen
	@echo "Generating mocks for authentication service..."
	@go generate ./...

install-mockgen:
	@bash ../../scripts/install-mockgen.sh

run-unit-tests: generate-mocks
	@echo "Running unit tests for authentication service..."
	@go test -v -tags=unit ./...
