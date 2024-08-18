# syntax = docker/dockerfile:1.3.0

FROM node:20.5.1-alpine AS build

WORKDIR /app/client

RUN apk add --update --no-cache openjdk8-jre-base

COPY ./client/package.json ./client/package-lock.json ./
RUN --mount=type=cache,target=/usr/src/app/.npm \
  npm set cache /usr/src/app/.npm && \
  npm install

COPY ./client/scripts ./scripts
COPY ./docs /app/docs
RUN npm run gen-api

COPY ./client/ ./
RUN npm run build

FROM caddy:2.8.4-alpine

COPY --from=build /app/client/dist/ /usr/share/caddy/
COPY ./deploy/Caddyfile /etc/caddy/Caddyfile

ENTRYPOINT ["caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
