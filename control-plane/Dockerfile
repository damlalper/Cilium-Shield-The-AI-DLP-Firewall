
# Stage 1: Build the Go application
FROM golang:1.21-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application, creating a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server .

# Stage 2: Create the final, minimal image
FROM scratch

# Copy the static binary from the builder stage
COPY --from=builder /server /server

# Expose the port the server listens on
EXPOSE 8080

# Set the entrypoint for the container
CMD ["/server"]
