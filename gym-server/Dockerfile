FROM golang:alpine AS builder

WORKDIR /gym-server
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/gymserver/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o migrate ./cmd/migrator/migrator.go

# Final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /gym-server/main .
COPY --from=builder /gym-server/migrate .

COPY config ./config
COPY migrations ./migrations

COPY entrypoint.sh ./
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]