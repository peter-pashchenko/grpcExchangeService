FROM golang:1.22-alpine AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download && go mod verify

COPY ./db .

RUN go build  -o ./app ./cmd/service

FROM alpine:3.14

RUN mkdir /app
WORKDIR /app

COPY --from=builder  /app/app .
COPY --from=builder  /app/db db


CMD ["./app" ]
