FROM golang:1.24

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Expose the application port
EXPOSE 8000

# Build and run the application
CMD go build -o backend ./cmd/api/main.go && ./backend