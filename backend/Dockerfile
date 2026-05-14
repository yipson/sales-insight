# Backend Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalar dependencias de build
RUN apk add --no-cache git

# Copiar go.mod y go.sum primero (cache de capas Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar binario
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Imagen final ligera
FROM alpine:latest

WORKDIR /app

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Copiar binario desde builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
