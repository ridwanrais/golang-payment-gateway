# Use an official Go runtime as the base image
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application. Replace "main.go" with the appropriate entry point to your application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o payment-gateway main.go

# Use a scratch (empty) image for a smaller final image
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/payment-gateway /app/payment-gateway

# Expose any necessary ports. Replace "8080" with your application's port if different.
EXPOSE 8090

# Command to run when the container starts
CMD ["/app/payment-gateway"]
