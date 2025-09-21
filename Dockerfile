# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23

ADD . /usr/src/app

WORKDIR /usr/src/app

CMD ["go", "run", "cmd/main.go"]