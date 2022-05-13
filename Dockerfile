FROM golang:1.17 as builder
WORKDIR /db-server
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /db-server main.go
# Финальный этап, копируем собранное приложение
FROM debian:buster-slim
COPY --from=builder /db-server/main /main
ENTRYPOINT ["/main", "-v"]