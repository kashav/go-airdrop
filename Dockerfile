FROM iron/go:dev

WORKDIR /app

ENV SRC_DIR=/go/src/github.com/kshvmdn/rdrp

ADD . $SRC_DIR

RUN cd $SRC_DIR && \
    go get && \
    go build -o rdrp && \
    cp rdrp /app/

ENTRYPOINT ["./rdrp"]
