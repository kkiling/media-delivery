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
FROM lscr.io/linuxserver/baseimage-ubuntu:jammy

# Установка runtime зависимостей
RUN apt-get update
RUN apt-get install -y mkvtoolnix
RUN apt-get install -y locales \
    && locale-gen en_US.UTF-8 \
    && update-locale LANG=en_US.UTF-8
RUN rm -rf /var/lib/apt/lists/

# Копируем бинарник из этапа сборки
COPY --from=builder /server_app /app/server_app
COPY --from=builder /app/swagger /swagger/

# Настройка переменных окружения
ENV CONFIG_PATH=/config/config.yml
ENV LANG=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8

# Экспортируем порты
EXPOSE 8080 8181

# Запуск через s6-setuidgid abc — тогда LSIO уже настроит UID/GID
CMD ["env", "LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8", "s6-setuidgid", "abc", "/app/server_app", "-c", "/config/config.yaml"]