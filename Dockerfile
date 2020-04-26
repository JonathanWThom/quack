FROM golang:1.14-alpine AS build_env

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /usr/src/app/quack .
FROM alpine 

COPY --from=build_env /usr/src/app/quack /bin/quack

