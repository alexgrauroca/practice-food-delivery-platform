generate-mocks: install-mockgen
	@echo "Generating mocks for pkg..."
	@go generate ./...

install-mockgen:
	@bash ../scripts/install-mockgen.sh
	@go get go.uber.org/mock/mockgen