# Use the official Go image as the build environment
FROM golang:1.22 AS builder

# Set Work Directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download Dependencies
RUN go mod download

# Copy Source Code
COPY . .

# Build App
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Using alpine as runtime
FROM alpine:latest  

# Copy the built executable file from the builder stage
COPY --from=builder /app/main /main

# Run App
CMD ["/main"]
