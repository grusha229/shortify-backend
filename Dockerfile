# Шаг 1: Сборка приложения
FROM golang:1.22.0 AS build

WORKDIR /app

COPY . .

# Копируем .env файл
COPY dev.env .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Шаг 2: Создаем минимальный образ для запуска
FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .
COPY --from=build /app/dev.env .  

EXPOSE 8080
CMD ["./main"]
