# Output binary name
BINARY_NAME = americano

# Package path for the main package
MAIN_PKG = ./cmd/americano

# Default target: Build and run the application
all: build run

# Build the Go project
build:
	go build -o ${BINARY_NAME} ${MAIN_PKG}

# Run the built binary
run: build
	./$(BINARY_NAME)

# Clean the build files
clean:
	rm -f $(BINARY_NAME)

# Run the application without building a binary
run-dev:
	go run $(MAIN_PKG)
