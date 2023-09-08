# Use an official Golang runtime as a parent image
FROM golang:1.16 AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go application
RUN go build -o myapp

# Start with a fresh Alpine Linux image
FROM alpine:latest

# Set the working directory to /app
WORKDIR /app

# Copy the compiled binary from the builder image
COPY --from=builder /app/myapp /app/myapp

# Expose the port your application will run on
EXPOSE 8080

# Define environment variables (if needed)
# ENV PORT4=8080
# ENV HOSTSITE=myhostsite

# Run your application
CMD ["./myapp"]
