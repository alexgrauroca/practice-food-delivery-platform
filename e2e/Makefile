run-e2e-tests: start-docker
	@echo "Running end-to-end tests..."
	@go test -v ./...

start-docker:
	@echo "Starting Docker containers for end-to-end tests..."
	@docker compose -f ../docker-compose.yml up -d --build