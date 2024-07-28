# Use a base image that has the required GLIBC versions
FROM golang:1.22.5 AS builder

WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download and cache the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go_server .

# Use a minimal image to run the binary
FROM debian:buster-slim

WORKDIR /root/

# Copy the built binary from the builder
COPY --from=builder /app/go_server .

# Expose the port the application runs on
EXPOSE 8080

# Run the binary
CMD ["./go_server"]
