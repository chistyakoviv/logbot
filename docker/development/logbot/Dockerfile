FROM golang:1.25.4-alpine3.22

# Allow to run programs with -race flag
ENV CGO_ENABLED=1

# build-base allows to run programs with -race flag
RUN apk add --no-cache build-base entr git openssl make

# Change working directory after installing dependencies
# (in case some deps are needed to be dowlaoded, unpacked and copied to the right place)
WORKDIR /go/src/github.com/chistyakoviv/logbot

RUN go install github.com/air-verse/air@latest

# Do not trigger reinstalling packages after updating dependencies
COPY ./go.mod ./go.sum ./

# Download dependencies required by go.mod
RUN go mod download

# Error: JSON arguments recommended for ENTRYPOINT/CMD to prevent unintended behavior related to OS signals
# Description:
# ENTRYPOINT and CMD instructions both support two different syntaxes for arguments:
# Shell form: CMD my-cmd start
# Exec form: CMD ["my-cmd", "start"]
# When you use shell form, the executable runs as a child process to a shell, which doesn't pass signals.
# This means that the program running in the container can't detect OS signals like SIGTERM and SIGKILL and respond to them correctly.
# see https://docs.docker.com/reference/build-checks/json-args-recommended/
CMD ["air", "-c", ".air.toml"]