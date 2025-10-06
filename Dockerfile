# Dockerfile

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install go-licenses and swag tools
RUN go install github.com/google/go-licenses@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs, gather licenses, and build the application
RUN /go/bin/swag init -g cmd/server/main.go -o ./docs
RUN /go/bin/go-licenses save ./... --save_path=/app/licenses
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/licenses ./licenses
COPY --from=builder /app/docs ./docs
COPY configs/config.yaml ./configs/config.yaml

EXPOSE 8080

CMD ["./main"]
