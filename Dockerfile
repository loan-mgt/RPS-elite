# Start from the official golang image
FROM golang:1.22-alpine AS builder

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=$TARGETPLATFORM

# Move to working directory /build
WORKDIR /build

# Copy and download dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the code into the container
COPY . .

# Move to the directory where main.go is located
WORKDIR /build/cmd

# Build the application
RUN go build -o /app .

# Start a new stage from scratch
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app /app

# Copy the necessary files for the application
COPY --from=builder /build/internal /build/internal
COPY --from=builder /build/static /static

# Command to run the executable
CMD ["/app"]
