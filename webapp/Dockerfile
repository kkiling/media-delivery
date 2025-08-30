# Stage 1: Build
FROM node:18-alpine as builder

WORKDIR /app

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm ci

# Копируем исходный код
COPY . .

# Собираем приложение
RUN npm run build

# Stage 2: Production
FROM nginx:alpine

# Копируем собранное приложение
COPY --from=builder /app/build /usr/share/nginx/html

# Копируем конфигурацию nginx
COPY nginx.conf /etc/nginx/templates/default.conf.template

# Устанавливаем переменную по умолчанию
ENV API_URL=http://localhost:8083

# Экспонируем порт
EXPOSE 80

# Запускаем nginx с поддержкой env переменных
CMD ["nginx", "-g", "daemon off;"]