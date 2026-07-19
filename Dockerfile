# --- STAGE 1: Build the application ---
FROM golang:alpine AS builder

# Install system dependencies (needed for templ or CGO if used)
RUN apk add --no-cache git npm nodejs

# Install templ CLI globally in the builder stage
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# Copy dependency manifests first for caching layers
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Generate templ components (converts .templ to .go files)
RUN templ generate

# Build the Go binary (disable CGO for a fully static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- STAGE 2: Final lightweight runner ---
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy any static assets or migration files if your app needs them at runtime
# COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./main"]
