FROM golang:alpine as builder

ENV GO111MODULE=on

RUN apk add -U git && go install github.com/swaggo/swag/cmd/swag@v1.6.9
WORKDIR /app
COPY . /app
RUN go get -d -v ./... && swag init

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o /gomod/papelito .

FROM node:14.20.1-alpine3.16 as node_builder

WORKDIR /app
COPY package.json /app
COPY ./.parcelrc /app/.parcelrc
COPY ./yarn.lock /app
COPY ./frontend /app/frontend

RUN yarn && \
    yarn build

FROM alpine:3.12

RUN apt update; \
	apk add openssl ca-certificates procps psmisc vim; \
	rm -rf /var/cache/apk


WORKDIR /app

COPY ./templates /app/templates
COPY CHANGELOG.MD /app/CHANGELOG.MD

COPY --from=builder /gomod/papelito /app/papelito
COPY --from=node_builder /app/assets /app/assets

RUN sed -e "s|http://localhost:9090|{{.ServerURL}}|g" /app/assets/index.html > /app/templates/layout/base.html

RUN chmod +x /app/papelito

CMD ["/app/papelito"]


