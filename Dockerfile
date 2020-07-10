FROM golang:1.14

WORKDIR /app

ADD . /app

RUN go build -o rdrp -v ./cmd/rdrp

ENTRYPOINT ["./rdrp"]
