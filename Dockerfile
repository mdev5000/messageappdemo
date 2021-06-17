FROM golang:1.16-alpine

ADD . /go/src/messageappdemo

WORKDIR /go/src/messageappdemo

RUN go install ./cmd/messageappdemo/messageappdemo.go

CMD ["messageappdemo"]
