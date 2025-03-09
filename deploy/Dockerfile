# syntax = docker/dockerfile:1

FROM node:23.9.0-alpine AS client

WORKDIR /app/client

RUN apk add --update --no-cache openjdk11-jre-headless

COPY ./client/package.json ./client/package-lock.json ./
RUN --mount=type=cache,target=/usr/src/app/.npm \
  npm set cache /usr/src/app/.npm && \
  npm install

COPY ./client/scripts ./scripts
COPY ./docs /app/docs
RUN npm run gen-api

COPY ./client/ ./
RUN npm run build

FROM golang:1.24.1-alpine AS server

RUN apk add --update --no-cache git alpine-sdk

WORKDIR /go/src/github.com/mazrean/one-poll/server

COPY ./server/go.mod ./server/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go mod download

COPY ./server .
COPY --from=client /app/client/dist/ ./handler/v1/static/
RUN --mount=type=cache,target=/root/.cache/go-build \
  go build -o one-poll -ldflags "-s -w"

FROM alpine:3.21

WORKDIR /go/src/github.com/mazrean/one-poll

RUN apk --update --no-cache add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk del tzdata \
  && mkdir -p /usr/share/zoneinfo/Asia \
  && ln -s /etc/localtime /usr/share/zoneinfo/Asia/Tokyo
RUN apk --update --no-cache add ca-certificates \
  && update-ca-certificates \
  && rm -rf /usr/share/ca-certificates

COPY --from=server /go/src/github.com/mazrean/one-poll/server/one-poll ./one-poll

ENTRYPOINT ["./one-poll"]
