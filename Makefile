.PHONY: build clean install lint

all: build

build:
	go build -v -o ./rdrp

clean:
	rm -f ./rdrp

install:
	go install

lint:
	$(GOPATH)/bin/golint .
