# Stage 1: Build the Go application
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/frieren-api ./cmd/api

# Stage 2: Run the built application
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/frieren-api .

EXPOSE 8000

CMD ["./frieren-api"]
