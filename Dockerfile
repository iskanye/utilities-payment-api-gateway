# Сборка
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/api-gateway ./cmd/api-gateway/main.go

# Запуск
FROM alpine
USER root
WORKDIR /home/app
COPY --from=builder /bin/api-gateway ./
COPY --from=builder /app/config ./config
ENTRYPOINT ["./api-gateway"]
CMD ["-config", "./config/dev.yaml"]