FROM golang:1.20.6-alpine3.18 AS builder

COPY . /github.com/IrinaFosteeva/TelegramBotPocket/
WORKDIR /github.com/IrinaFosteeva/TelegramBotPocket/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/IrinaFosteeva/TelegramBotPocket/bin/bot .
COPY --from=0 /github.com/IrinaFosteeva/TelegramBotPocket/configs configs/

EXPOSE 80

CMD ["./bot"]