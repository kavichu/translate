FROM golang:1.8.1-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git

ENV GOPATH /go

RUN go get -u cloud.google.com/go/translate

ADD main.go /go/

RUN go build -o main main.go

ENTRYPOINT ["/go/main"]