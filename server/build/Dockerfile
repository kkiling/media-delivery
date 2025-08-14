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
RUN rm -rf /var/lib/apt/lists/

# Копируем бинарник из этапа сборки
COPY --from=builder /server_app /app/server_app
COPY --from=builder /app/swagger /swagger/

# Настройка переменных окружения
ENV CONFIG_PATH=/config/config.yml

# Экспортируем порты
EXPOSE 8080 8181

# Переключаемся на пользователя abc (учитывает PUID/PGID)
#USER abc

# Запуск от имени пользователя abc (автоматически учитывает PUID/PGID)
#CMD ["/server_app", "-c", "/config/config.yaml"]

# Запуск через s6-overlay (LSIO сам запустит от abc с PUID/PGID)
#CMD ["/app/server_app", "-c", "/config/config.yaml"]
# Запуск через s6-setuidgid abc — тогда LSIO уже настроит UID/GID
CMD ["s6-setuidgid", "abc", "/app/server_app", "-c", "/config/config.yaml"]
