FROM golang:1.11

WORKDIR /go/src/wrg/rts/lab/
COPY . .

RUN go get ./...

RUN go build main.go

CMD ["./main"]