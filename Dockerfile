# Этап сборки
FROM golang:1.24 AS build

WORKDIR /src

# Копируем и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение со статической линковкой
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app ./app/server

# Финальный этап с Distroless
FROM gcr.io/distroless/static-debian11:nonroot

# Создаем директорию для приложения
WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=build /go/bin/app ./app
COPY --from=build /src/config.yaml ./config.yaml

# Определяем порт
EXPOSE 8080

# Запуск от nonroot пользователя
USER nonroot:nonroot

# Запускаем приложение
ENTRYPOINT ["/app/app"]
