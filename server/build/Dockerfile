# Этап сборки (builder)
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o /server_app ./cmd/main.go

# Финальный этап (runtime)
FROM lscr.io/linuxserver/code-server:latest

# Установка runtime зависимостей
RUN apt-get update
RUN apt-get install -y mkvtoolnix
RUN rm -rf /var/lib/apt/lists/

# Копируем бинарник из этапа сборки
COPY --from=builder /server_app /app/server_app

# Настройка переменных окружения
ENV CONFIG_PATH=/config/config.yml

# Рабочая директория
WORKDIR /app

# Экспортируем порты
EXPOSE 8080 8181

# Команда запуска (используем shell форму для подстановки переменных)
CMD ["sh", "-c", "/app/server_app -c \"${CONFIG_PATH}\""]