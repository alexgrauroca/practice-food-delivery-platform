generate-mocks: install-mockgen
	@echo "Generating mocks for authentication service..."
	@go generate ./...

install-mockgen:
	@bash ../../scripts/install-mockgen.sh

run-integration-tests: start-docker
	@echo "Running integration tests for authentication service..."
	@go test -v -tags=integration ./... -run IntegrationTest

start-docker:
	@echo "Starting Docker containers for integration tests..."
	@docker-compose -f ../../docker-compose.yml up -d
