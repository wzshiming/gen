FROM golang:alpine AS builder
WORKDIR /go/src/github.com/wzshiming/gen
COPY . .
RUN apk add -U --no-cache git
RUN go get ./cmd/gen
CMD gen
