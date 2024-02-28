# Use a lightweight base image
FROM golang:alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY ./src/ .

# Build the Go application
RUN go build -o app

# Use a minimal base image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app /app/app

# Make the binary executable
RUN chmod +x /app/app

# Run the application
CMD ["/app/app"]

