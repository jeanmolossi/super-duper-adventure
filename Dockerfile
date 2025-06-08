FROM golang:1.24-alpine

RUN apk add --update bash curl

WORKDIR /go/src

ENTRYPOINT [ "top" ]