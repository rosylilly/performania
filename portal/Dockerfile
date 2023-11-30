# syntax=docker/dockerfile:1

ARG NODE_VERSION=20.9.0
ARG GO_VERSION=1.21.4

FROM node:${NODE_VERSION}-slim as node

WORKDIR /app/www
COPY www /app/www

RUN yarn install --frozen-lockfile
RUN yarn build

FROM golang:${GO_VERSION} as go

ENV GO111MODULE=on
ENV BUILD=release

WORKDIR /app
COPY . /app

RUN go mod download

COPY --from=node /app/www/dist /app/www/dist

RUN make build

FROM scratch

EXPOSE 8080

COPY --from=go /app/bin/portal /bin/portal
CMD [ "/bin/portal" ]
