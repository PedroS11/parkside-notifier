# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.12-alpine base image
FROM golang:1.24

# Add Maintainer Info
LABEL maintainer="Pedro Silva <pedrosilva1137work@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./src .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]
