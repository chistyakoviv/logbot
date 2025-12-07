FROM golang:1.25.4-alpine3.22 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/logbot ./cmd/logbot

FROM alpine:3.22 AS runner

RUN apk add --no-cache curl
WORKDIR /app

COPY --from=builder /app/bin/logbot ./

EXPOSE 80

HEALTHCHECK --interval=5s --timeout=5s --start-period=1s CMD curl --fail http://127.0.0.1:80/healthcheck || exit 1

CMD ["./logbot"]
