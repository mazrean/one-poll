# syntax = docker/dockerfile:1.3.0

FROM golang:1.17.5-alpine AS build

RUN apk add --update --no-cache git alpine-sdk

WORKDIR /go/src/github.com/mazrean/one-poll/server

COPY ./server/go.mod ./server/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go mod download

COPY ./server .
COPY ./docs /go/src/github.com/mazrean/one-poll/docs
RUN go generate ./...
RUN --mount=type=cache,target=/root/.cache/go-build \
  go build -o one-poll -ldflags "-s -w"

FROM alpine:3.15.0

WORKDIR /go/src/github.com/mazrean/one-poll

RUN apk --update --no-cache add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk del tzdata \
  && mkdir -p /usr/share/zoneinfo/Asia \
  && ln -s /etc/localtime /usr/share/zoneinfo/Asia/Tokyo
RUN apk --update --no-cache add ca-certificates \
  && update-ca-certificates \
  && rm -rf /usr/share/ca-certificates

COPY --from=build /go/src/github.com/mazrean/one-poll/server/one-poll ./one-poll

ENTRYPOINT ["./one-poll"]
