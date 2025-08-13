# Этап сборки (добавлен AS builder)
FROM golang:1.24 AS builder

# Установка runtime зависимостей (исправлено для Alpine)
RUN apt update
RUN apt install -y mkvtoolnix

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o /server_app ./cmd/main.go

# Переменные окружения
ENV CONFIG_PATH=configs/config.yml

# Экспортируем порт
EXPOSE 8080

# Команда запуска
CMD /server_app -c "${CONFIG_PATH}"