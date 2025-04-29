FROM golang:alpine

WORKDIR /gym-server
COPY ./ /gym-server

RUN go mod download
RUN go install