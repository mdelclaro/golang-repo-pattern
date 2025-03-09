# Test the application
test: 
	gotestsum --format pkgname ./...

# Build the application
build-dependencies:
	docker compose build
	
# Run the application
run:
	docker compose up

# Generate mocks
mock:
	mockgen -package=mock -source=./internal/domain/device/repository.go -destination=internal/domain/device/mock/repository.go
	