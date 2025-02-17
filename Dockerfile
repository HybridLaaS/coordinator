# Build stage
FROM golang:1.23 AS builder

WORKDIR /build
COPY . .

# Ensure dependencies are resolved
RUN go mod tidy

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o coordinator

# Minimal runtime image
FROM alpine:latest

WORKDIR /coordinator

# Copy the built binary
COPY --from=builder /build/coordinator /coordinator/coordinator
COPY ./.env /coordinator/.env

# Expose the application port
EXPOSE 8090

# Run the binary
CMD ["/coordinator/coordinator"]