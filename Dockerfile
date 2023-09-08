# Use the desired Go version
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app

# Copy only the necessary Go module files
COPY go.mod go.sum ./

# Download and cache Go module dependencies
RUN go mod download

# Copy the rest of the project into the container
COPY . .

# Build the Go application with explicit flags
RUN go build -o myapp

# Use a minimal base image for the final container
FROM alpine:latest

# Set the working directory in the final container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the port your Go application listens on (e.g., 8080)
EXPOSE 8080

# Define environment variables (if needed)
ENV PORT3=8080
ENV HOSTSITE=https://resoledge.com

# Run your Go application
CMD ["./myapp"]
