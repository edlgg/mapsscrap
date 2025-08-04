BINARY_NAME := mapsscrap

# Default binary name
BINARY_NAME=mapsscrap

# Build the application
build:
	go build -o bin/$(BINARY_NAME) main.go

# Run the application with sample parameters for Mexico City
run:
	go run main.go --lat 19.4343491 --lon -99.1775742 --query "dentist" --radius 2

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f $(BINARY_NAME)

# Tidy and verify dependencies
tidy:
	go mod tidy
	go mod verify

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go get -v ./...

# Format code
fmt:
	go fmt ./...