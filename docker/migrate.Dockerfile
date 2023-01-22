# Initial stage: download modules
FROM golang:1.19-alpine as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download
RUN go build ./cmd/migration/main.go