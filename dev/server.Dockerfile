# syntax = docker/dockerfile:1

FROM golang:1.23.3-alpine AS build

RUN apk --update --no-cache add tzdata && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
  apk del tzdata

RUN apk add --update --no-cache git

RUN --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/cosmtrek/air@v1.27.3

RUN --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/go-delve/delve/cmd/dlv@v1.7.2

WORKDIR /go/src/github.com/mazrean/one-poll

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go mod download

ENTRYPOINT ["air"]
CMD ["-c", ".air.toml"]
