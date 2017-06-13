FROM iron/go:dev

WORKDIR /app

ENV SRC_DIR=/go/src/github.com/kshvmdn/rdrp

ADD . $SRC_DIR

RUN cd $SRC_DIR && \
    go get -u -v ./... && \
    go build -o rdrp -v ./cmd/rdrp && \
    cp rdrp /app/

ENTRYPOINT ["./rdrp"]
