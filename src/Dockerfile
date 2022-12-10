FROM golang:1.16-alpine

ADD . /api
WORKDIR /api

RUN go mod download

RUN go build -o /audioApi

# RUN mkdir /data

EXPOSE 8080

CMD [ "/audioApi" ]