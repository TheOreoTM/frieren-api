# Use an official lightweight Go image as the build stage
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application using the Makefile
RUN make build

# Use a minimal image for running the application
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/frieren-api .

# Expose the application's port
EXPOSE 8000

# Run the application
CMD ["./frieren-api"]
