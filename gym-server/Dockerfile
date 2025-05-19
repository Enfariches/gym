FROM golang:alpine

WORKDIR /gym

RUN apk add --no-cache git bash

# Скопируем только необходимые ресурсы (остальное монтируется с volume)
COPY config ./config
COPY migrations ./migrations
COPY lib/pdf ./lib/pdf

COPY entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]