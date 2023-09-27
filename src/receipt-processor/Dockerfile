# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and sum files to enable efficient Go module caching
COPY go.mod go.sum ./

# Download and install Go dependencies (this step is cached if the go.mod and go.sum files don't change)
RUN go mod download

# Copy the rest of the local package files to the container's workspace
COPY . .

# Build the Go application inside the container
RUN go build -o main

# Expose a port (assuming your Go app listens on port 8080)
EXPOSE 8080

# Run the Go application when the container starts
CMD ["./main"]
