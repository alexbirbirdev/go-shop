# Сборка
FROM golang:1.23 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/go-shop ./cmd/

# Запуск
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/go-shop .
COPY --from=builder /app/config ./config
COPY --from=builder /app/.env .

EXPOSE 8080
CMD [ "./go-shop" ]
