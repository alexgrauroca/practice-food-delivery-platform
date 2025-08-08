openapi-validate-sync:
	@cd ../.. && ./scripts/openapi-validate-sync.sh services/customer-service

openapi-validate-docs:
	@cd ../.. && ./scripts/openapi-validate-docs.sh services/customer-service

openapi-generate:
	@cd ../.. && ./scripts/openapi-bundle.sh services/customer-service

generate-clients: install-openapi-generator
	@echo "Generating clients for customer service..."
	@cd ../.. && mkdir -p clients/customer-service
	@npx @openapitools/openapi-generator-cli generate -i docs/dist/openapi.yaml \
    		-g go \
    		-o ../../clients/customer-service \
    		--package-name customerclient
	@cd ../../clients/customer-service && \
		if [ ! -f go.mod ]; then \
			go mod init github.com/alexgrauroca/practice-food-delivery-platform/clients/customer-service; \
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
