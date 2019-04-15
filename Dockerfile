FROM golang:1.11

WORKDIR /svr/
COPY . .

RUN go build main.go

CMD ["./main"]