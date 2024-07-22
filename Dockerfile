# Use the official Golang image to create a build artifact.
FROM golang:1.21.3-alpine3.18 AS builder

# Set the Current Working Directory inside the container
WORKDIR /media-service/

COPY ./../media-service.bin .

# Copy go mod and sum files from the parent directory
COPY go.mod go.sum ./

# Copy other necessary files from the parent directory
COPY .env .

# Declare the volume to persist image storage
VOLUME ["/media-service/image-storage"]

RUN chmod +x media-service.bin
# This container exposes port 8082 to the outside world
EXPOSE 8082

# Run the binary program
CMD ["./media-service.bin"]