FROM golang:alpine AS builder
RUN apk add --no-cache gcc git musl-dev make zip
RUN mkdir /app

WORKDIR /app
COPY . .
RUN go mod download

ARG PLUGIN
WORKDIR /app/docker/plugin
RUN go mod download

RUN go build -buildmode=plugin -o $PLUGIN.so $PLUGIN.go
