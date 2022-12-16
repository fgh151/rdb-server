FROM golang:1.18 as builder
WORKDIR /db-server
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN go build -o /db-server main.go && go build -o /db-server/cli cli.go && go build  -buildmode=plugin -o smsc.so /db-server/plugins/smsc/smsc_plugin.go

# Финальный этап, копируем собранное приложение
FROM debian:buster-slim
COPY --from=builder /db-server/main /main
COPY --from=builder /db-server/cli /cli
COPY --from=builder /db-server/smsc.so /smsc.so
ENTRYPOINT ["/main", "-v"]
