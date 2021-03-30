FROM golang:1.15-alpine3.12 AS builder

RUN go version

COPY . /reminder_bot/
WORKDIR /reminder_bot/

RUN go mod download
RUN GOOS=linux go build -o tg_bot_service cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /reminder_bot/tg_bot_service .
COPY --from=0 /reminder_bot/config /config

CMD ["./parse_service"]