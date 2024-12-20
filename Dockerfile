# Use the official Go image as a base image for building
FROM golang:1.21.13 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Install goose for database migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build the application
RUN go build -o Flux

# Use a smaller base image for the final container
FROM debian:latest

# Set working directory
WORKDIR /app

# Copy the built binary and necessary files from the builder stage
COPY --from=builder /app/Flux .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY sql /app/sql

# Expose the port your app listens on
EXPOSE 8080

# Command to run database migrations and start the app
CMD ["sh", "-c", "goose -dir=/app/sql/schema postgres $DB_URL up && ./Flux"]
