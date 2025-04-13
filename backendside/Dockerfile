# Use the official Golang image as the base image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY backend .

# Build the application
RUN go build -o main .

# Use a smaller base image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
