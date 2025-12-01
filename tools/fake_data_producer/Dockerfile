# Dockerfile для fake_data_producer
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Сначала зависимости
COPY go.mod go.sum ./
RUN go mod download

# Затем весь исходный код (каждое изменение кода будет инвалидировать слой)
COPY . .

# Явно собираем бинарь под linux без CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /fake-producer ./tools/fake_data_producer

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /fake-producer ./fake-producer

# Устанавливаем средства для ввода
RUN apk add --no-cache bash

ENTRYPOINT ["./fake-producer"]