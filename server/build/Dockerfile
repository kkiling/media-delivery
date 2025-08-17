# Этап сборки (builder)
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o /server_app ./cmd

# Финальный этап (runtime)
FROM lscr.io/linuxserver/baseimage-ubuntu:jammy

# Установка runtime зависимостей
RUN apt-get update
RUN apt-get install -y mkvtoolnix
RUN rm -rf /var/lib/apt/lists/

# Копируем бинарник из этапа сборки
COPY --from=builder /server_app /app/server_app
# Создаем сервис для s6-overlay
RUN mkdir -p /etc/services.d/server_app
COPY --from=builder /app/swagger /etc/services.d/server_app/swagger/
COPY --from=builder /app/migrations /etc/services.d/server_app/migrations/

# Экспортируем порты
EXPOSE 8080 8181

COPY <<EOF /etc/services.d/server_app/run
#!/bin/sh
export LANG="en_US.UTF-8"
export LC_ALL="en_US.UTF-8"
umask 0002
# запускаем от имени PUID/PGID (user=abc)
exec s6-setuidgid abc /app/server_app -c /configs/config.yaml
EOF

RUN chmod +x /etc/services.d/server_app/run

# s6-overlay автоматически запустит наш сервис