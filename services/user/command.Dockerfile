# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-buster AS build

WORKDIR /app


COPY . ./
RUN go mod download

RUN ls -lisa

RUN CGO_ENABLED=0 go build -o /command command/main.go

##
## Deploy
##
FROM alpine:3.15.0

WORKDIR /

COPY --from=build /command /user-command

ENTRYPOINT ["/user-command"]