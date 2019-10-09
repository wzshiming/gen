FROM golang:alpine AS builder
WORKDIR /go/src/github.com/wzshiming/gen
COPY . .
RUN apk add -U --no-cache git
RUN go install ./cmd/gen
CMD gen
