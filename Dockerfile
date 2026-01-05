# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/front

# Копируем package файлы и устанавливаем зависимости
COPY src/front/package*.json ./
RUN npm ci

# Копируем весь фронтенд код
COPY src/front/ ./

# Билдим фронтенд с переменной окружения (он соберется в ./dist)
RUN DOCKER_BUILD=true npm run build

# Stage 2: Build backend
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY src/back/go.mod src/back/go.sum ./
RUN go mod download

# Копируем весь backend код
COPY src/back/ ./

# Копируем собранный фронтенд из первого stage (теперь из правильного места!)
COPY --from=frontend-builder /app/front/dist ./web

# Билдим Go приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /laritmo ./cmd/server

# Stage 3: Final lightweight image
FROM alpine:latest

WORKDIR /app

# Устанавливаем CA сертификаты для HTTPS
RUN apk --no-cache add ca-certificates

# Копируем бинарник
COPY --from=backend-builder /laritmo /app/laritmo

# Копируем необходимые файлы
COPY --from=backend-builder /app/web ./web
COPY --from=backend-builder /app/configs ./configs
COPY --from=backend-builder /app/migrations ./migrations

# Создаем непривилегированного пользователя
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8443

CMD ["/app/laritmo"]