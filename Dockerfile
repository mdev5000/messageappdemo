FROM golang:1.16-alpine

ADD . /go/src/messageapidemo

WORKDIR /go/src/messageapidemo

RUN go install ./cmd/messageapidemo/messageapidemo.go

CMD ["messageapidemo"]
