
# 8. Clients Auto-Generator

## Status

Accepted

Date: 2025-08-12

## Context

As our service architecture grows, we need a standardized and automated way to generate client libraries for 
inter-service communication. Currently, services like the customer-service need to communicate with other services 
(e.g., authentication-service) through HTTP clients. Manual implementation of these clients is error-prone and 
time-consuming, especially when API changes occur.

The OpenAPI specifications already define our service interfaces, making them an ideal source for generating 
consistent client libraries.

## Decision

We will implement an automated client generation system with the following characteristics:

1. Use OpenAPI Generator to create Go client libraries from service OpenAPI specifications
2. Store all generated clients in a centralized `/clients` directory, with each service having its own subdirectory
3. Implement a standardized configuration for client generation using YAML files
4. Use Go modules replacement to manage client dependencies between services
5. Adapt build processes (Dockerfile and docker-compose) to support the new client structure
6. Adapt the CI workflow to check if the clients are up to date with the OpenAPI specifications
7. Implement a Makefile to automate client generation and dependency management

## Consequences

### Positive

- Automated and consistent client generation across all services
- Reduced development time and fewer manual errors in client implementations
- Type-safe client interfaces that automatically update with API changes
- Centralized management of client libraries in the `/clients` directory
- Better dependency management through Go modules replacement
- Simplified integration testing and service interaction

### Negative

- Additional build complexity in Docker configurations
- Need to maintain OpenAPI specifications strictly in sync with implementations
- Increased project size due to generated code
- Additional step in the build process to generate clients

### Neutral

- Services need to use replace directives in their go.mod files to reference local client packages
- Build context needs to include the entire project root for Docker builds
- Generated code requires regular updates when service APIs change

## Implementation Notes

1. **Client Generation Structure**:
   ```
   /clients
     /authentication-service
     /customer-service
     /[other-service]
   ```

2. **Client Configuration**:
    - Use `client-gen-config.yaml` for each service
    - Configure package names, output directories, and Git metadata
    - Enable Go submodule support for proper dependency management
    ```yaml
    # Example client generation configuration for authentication-service
    generatorName: go
    outputDir: ../../clients/authentication-service
    packageName: authclient
    gitUserId: alexgrauroca
    gitRepoId: practice-food-delivery-platform
    additionalProperties:
        module: github.com/alexgrauroca/practice-food-delivery-platform/authclient
        isGoSubmodule: "true"
        enumClassPrefix: "true"
        packageName: authclient
        disallowAdditionalPropertiesIfNotPresent: "true"

    ```

3. **Module Management**:
   ```bash
   # Example module replacement in service's go.mod
   replace github.com/alexgrauroca/practice-food-delivery-platform/authclient => ../../clients/authentication-service
   ```

4. **Docker Build Adaptation**:
    - Set build context to project root
    - Copy client libraries before building services
    - Ensure proper dependency resolution
    ```yaml
    # Example Docker Compose file for customer-service
    customer-service:
      build:
        context: . # Context needs to be set to the root of the project because of clients import
        dockerfile: ./services/customer-service/Dockerfile
      container_name: customer-service
    ```
   
    ```yaml
    # Example Dockerfile for customer-service
    # Build stage
    FROM golang:1.24-alpine AS builder

    # Copy clients
    WORKDIR /app/clients/authentication-service
    COPY ./clients/authentication-service/ ./

    # Copy and build the service
    WORKDIR /app/services/customer-service
    COPY ./services/customer-service/go.mod ./services/customer-service/go.sum ./
    RUN go mod download

    COPY ./services/customer-service/ ./

    RUN go build -o customer-service ./cmd/main.go

    # Final stage
    FROM alpine:latest

    COPY --from=builder /app/services/customer-service/customer-service .
    EXPOSE 8080

    CMD ["./customer-service"]
    ```

5. **Make Commands**:
   ```makefile
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
   ```

## Related Documents

- [OpenAPI Tools Generator](https://github.com/OpenAPITools/openapi-generator)

## Contributors

- Àlex Grau Roca