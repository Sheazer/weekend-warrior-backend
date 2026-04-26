# STAGE 1: Сборка с поддержкой CGO
FROM golang:1.25-alpine AS builder

# Устанавливаем gcc и musl-dev (нужны для компиляции SQLite)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ВКЛЮЧАЕМ CGO и собираем
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/server/main.go

# STAGE 2: Финальный образ
FROM alpine:latest

# Добавляем библиотеки для работы sqlite в рантайме
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируем бинарник из билдера
COPY --from=builder /app/main .
# Копируем файл базы (если он есть в корне проекта)
# COPY --from=builder /app/warriors.db .

EXPOSE 8080

CMD ["./main"]