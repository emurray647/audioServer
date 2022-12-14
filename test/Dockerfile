ARG GOVERSION=1.16
FROM golang:${GOVERSION} AS builder

ARG SRC_ROOT
ARG TEST_ROOT

RUN echo ${SRC_ROOT}
RUN echo ${TEST_ROOT}

ADD ${SRC_ROOT} /api
WORKDIR /api

ADD ${TEST_ROOT}/samples /samples

RUN go mod download

CMD go test -v --count=1 ./test/...