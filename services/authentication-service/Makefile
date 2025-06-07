run-unit-tests: generate-mocks
	@echo "Running unit tests for authentication service..."
	@go test -v ./...

generate-mocks: install-mockgen
	@echo "Generating mocks for authentication service..."
	@go generate ./...

install-mockgen:
	@bash ../../scripts/install-mockgen.sh

run-integration-tests: start-mongodb
	@echo "Running integration tests for authentication service..."
	@bash -c "set -a && source ../../.env && set +a && go test -v --tags=integration ./internal/customers/integration/..."

start-mongodb:
	@echo "Starting MongoDB Docker containers for integration tests..."
	@docker compose -f ../../docker-compose.yml up -d mongodb
