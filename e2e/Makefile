run-e2e-tests: start-docker
	@echo "Running end-to-end tests..."
	@go test -v -tags=e2e ./features/...

run-e2e-auth-tests: start-docker
	@echo "Running authentication end-to-end tests..."
	@go test -v -tags=authentication ./features/...

run-e2e-customer-tests: start-docker
	@echo "Running customer end-to-end tests..."
	@go test -v -tags=customer ./features/...

run-e2e-restaurant-tests: start-docker
	@echo "Running restaurant end-to-end tests..."
	@go test -v -tags=restaurant ./features/...

start-docker:
	@echo "Starting Docker containers for end-to-end tests..."
	@docker compose -f ../docker-compose.yml down
	@docker compose -f ../docker-compose.yml up -d --build