# Build stage

# Start from golang:1.12-alpine base image
FROM golang:1.24 AS builder

# Add Maintainer Info
LABEL maintainer="Pedro Silva <pedrosilva1137work@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./src ./src

# Build the Go app
RUN go build -o main ./src



# Run stage
FROM debian:bookworm-slim

# Install CA certificates for Telegram calls as debian:bookworm-slim image doesn't have them by default
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
