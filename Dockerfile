# Этап 1: Сборка
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum (для кэширования зависимостей)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники и собираем бинарник
COPY . .
RUN go build -o url-shortener .

# Этап 2: Финальный образ
FROM alpine:latest

# Создаём пользователя без прав
RUN adduser -D -g '' appuser

WORKDIR /app

# Копируем только бинарник
COPY --from=builder /app/url-shortener .

# Запускаем от имени непривилегированного пользователя
USER appuser

EXPOSE 8080
CMD ["./url-shortener"]
