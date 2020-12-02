FROM golang:1.15 AS build

RUN mkdir /app
COPY . /appp
WORKDIR /app

RUN CGO_ENABLED=0 go build -o pusher

FROM alpine:3.9
RUN apk add --update tzdata
RUN apk add --update ca-certificates
RUN apk add --update bash

COPY --from=build /app /
RUN chmod +x /pusher
