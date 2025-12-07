FROM alpine:3.22

# bash needed for wait-for-it
RUN apk add --no-cache make bash curl postgresql-client

RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh

COPY ./docker/common/wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod 555 /usr/local/bin/wait-for-it

RUN addgroup -g 1000 app && adduser -u 1000 -G app -s /bin/sh -D app

WORKDIR /app

COPY ./Makefile ./
COPY ./migrations ./migrations

USER app
