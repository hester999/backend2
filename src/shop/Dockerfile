FROM golang:1.24.3 AS builder


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go install github.com/swaggo/swag/cmd/swag@latest


RUN swag init --generalInfo cmd/main.go --dir ./

RUN make build

# Финальный образ
FROM ubuntu:22.04

WORKDIR /app

# Установим нужные зависимости для запуска Go-приложения
RUN apt update && apt install -y ca-certificates

# Копируем скомпилированный бинарник
COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/docs /app/docs

EXPOSE 8080

CMD ["/app/app"]
