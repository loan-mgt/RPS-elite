# Start from the official golang image
FROM golang:1.22-alpine AS builder

ARG GOARCH

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=$GOARCH

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o app .

# Start a new stage from scratch
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/app /app

# Copy the static directory
COPY  static /static

# Command to run the executable
CMD ["/app"]
