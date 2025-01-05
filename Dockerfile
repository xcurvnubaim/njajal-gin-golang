# Stage 1: Build
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
RUN go build -o main cmd/api/main.go

# Stage 2: Run
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port (change to your app's port)
EXPOSE 3000

# Command to run the application
CMD ["./main"]
