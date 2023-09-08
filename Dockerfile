# Use the desired Go version
FROM golang:1.16 AS builder

# Set the working directory
WORKDIR /app

# Copy the entire project into the container
COPY . .

# Clean the Go build cache
RUN go clean -cache

# Build the Go application with explicit flags
RUN CGO_ENABLED=1 GOARCH=amd64 go build -o myapp

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
