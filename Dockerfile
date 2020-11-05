ARG COMMAND=unknown
ARG GO_MODULE=unknown
FROM golang:1.14 as builder
ARG COMMAND
ARG GO_MODULE

RUN mkdir /opt/${COMMAND}
WORKDIR /opt/${COMMAND}

COPY . /opt/${COMMAND}

ENV GO111MODULE on
RUN go build ${GO_MODULE}/cmd/${COMMAND}

FROM debian:stretch-slim
ARG COMMAND

RUN apt-get update
RUN apt-get install -y ca-certificates

COPY --from=builder /opt/${COMMAND}/${COMMAND} /usr/local/bin/${COMMAND}

RUN ln -s /usr/local/bin/${COMMAND} /usr/local/bin/run_server
ENTRYPOINT /usr/local/bin/run_server
