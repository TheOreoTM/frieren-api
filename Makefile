# Define variables
BINARY_NAME=frieren-api
SRC=./cmd/api

# Default target: build the project
all: build

# Build the Go project
build:
	@echo "Building the project..."
	go build -o $(BINARY_NAME) $(SRC)

# Run the Go project
run: build
	@echo "Running the project..."
	./$(BINARY_NAME)

# Clean the generated files (binaries, CSVs, etc.)
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -f characters.csv

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt $(SRC)

# Run tests
test:
	@echo "Running tests..."
	go test $(SRC)

# Tidy up Go module dependencies
tidy:
	@echo "Tidying up dependencies..."
	go mod tidy

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download

# Specify that these are not actual file names
.PHONY: all build run clean fmt test tidy deps
