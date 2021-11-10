#FROM golang:1.14-alpine as builder
#RUN mkdir -p /go/src/github.com/alperhankendi/golang-api
#
#ENV GOPATH /go
#WORKDIR /go/src/github.com/alperhankendi/golang-api
#ADD go.mod .
#ADD go.sum .
#RUN go mod download
#ADD . /go/src/github.com/alperhankendi/golang-api
#
#RUN go build


FROM alpine

RUN mkdir -p /app
#COPY --from=builder /go/src/github.com/alperhankendi/golang-api /app/
WORKDIR /app
ADD golang-api /app
ADD config/config.dev.yaml config.dev.yaml
RUN chmod +x /app/golang-api

ENTRYPOINT ["golang-api"]