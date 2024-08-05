# Use a base image that contains Go 1.20
FROM golang:1.22.5 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a more comprehensive base image for the final stage
FROM debian:latest

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates

# Set the Working Directory inside the container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Run the executable
CMD ["./main"]
