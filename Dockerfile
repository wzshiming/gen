FROM golang:alpine3.8 AS builder
WORKDIR /go/src/github.com/wzshiming/gen
COPY . .
RUN go install ./cmd/...
CMD gen
