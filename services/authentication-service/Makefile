openapi-generate:
	@cd ../.. && ./scripts/openapi-bundle.sh services/authentication-service
	@$(MAKE) openapi-validate-docs

openapi-validate-docs:
	@cd ../.. && ./scripts/openapi-validate-docs.sh services/authentication-service

openapi-validate-sync:
	@cd ../.. && ./scripts/openapi-validate-sync.sh services/authentication-service

generate-clients: install-openapi-generator
	@echo "Generating clients for authentication service..."
	@cd ../.. && mkdir -p clients/authentication-service
	@npx @openapitools/openapi-generator-cli generate \
			-i docs/dist/openapi.yaml \
    		-c client-gen-config.yaml
	@cd ../../clients/authentication-service && \
		if [ ! -f go.mod ]; then \
			go mod init github.com/alexgrauroca/practice-food-delivery-platform/clients/authentication-service; \
		fi && \
		go mod tidy

install-openapi-generator:
	@bash ../../scripts/install-openapi-generator.sh

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
