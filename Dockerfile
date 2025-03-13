FROM golang:alpine

WORKDIR /gym-server
COPY ./ /gym-server

RUN go mod download
RUN go get .
RUN go build
RUN go install

ENTRYPOINT [ "health", "api" ]