run-e2e-tests: start-docker
	@echo "Running end-to-end tests..."
	@go test -v -tags=e2e ./...

run-e2e-auth-tests: start-docker
	@echo "Running authentication end-to-end tests..."
	@go test -v -tags=authentication ./...

run-e2e-customer-tests: start-docker
	@echo "Running customers end-to-end tests..."
	@go test -v -tags=customer ./...

start-docker:
	@echo "Starting Docker containers for end-to-end tests..."
	@docker compose -f ../docker-compose.yml up -d --build