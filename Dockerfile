# Development stage
FROM golang:1.23-alpine

WORKDIR /app

# Install git (required for fetching dependencies)
RUN apk add --no-cache git

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Run the Go application
CMD ["go", "run", "./cmd/transaction/main.go"]