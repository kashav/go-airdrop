NAME := rdrp
MAIN := ./cmd/$(NAME)
SRCS := $(shell find . -type f -name '*.go')
PKGS := $(shell go list ./...)

build = GOOS=$(1) GOARCH=$(2) go build -o build/$(NAME)$(3) $(MAIN)
tar = cd build && tar -cvzf $(1)_$(2).tar.gz $(NAME)$(3) && rm $(NAME)$(3)
zip = cd build && zip $(1)_$(2).zip $(NAME)$(3) && rm $(NAME)$(3)

.PHONY: all build clean fmt fmt-save get-tools lint test vet
.DEFAULT: all

all: rdrp

get-tools:
	@echo "+ $@"
	@go get -u -v golang.org/x/lint/golint

clean:
	@echo "+ $@"
	rm -rf build $(NAME)
	mkdir -p build

rdrp: $(SRCS)
	@echo "+ $@"
	@go build -o $(NAME) $(MAIN)

fmt:
	@echo "+ $@"
	@test -z "$$(gofmt -s -l . 2>&1)" || \
		(echo >&2 "+ please run 'gofmt -s' or 'make fmt-save'" && false)

fmt-save:
	@echo "+ $@"
	@gofmt -s -l . 2>&1 | xargs gofmt -s -l -w

vet:
	@echo "+ $@"
	@go vet $(PKGS)

lint:
	@echo "+ $@"
	@golint -set_exit_status=1 $(PKGS)

test:
	@echo "+ $@"
	@go test -race -v $(PKGS)

build: darwin linux windows

darwin: build/darwin_amd64.tar.gz

build/darwin_amd64.tar.gz: $(SRCS)
	$(call build,darwin,amd64,)
	$(call tar,darwin,amd64)

linux: build/linux_arm.tar.gz build/linux_arm64.tar.gz build/linux_386.tar.gz build/linux_amd64.tar.gz

build/linux_386.tar.gz: $(SRCS)
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz: $(SRCS)
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

build/linux_arm.tar.gz: $(SRCS)
	$(call build,linux,arm,)
	$(call tar,linux,arm)

build/linux_arm64.tar.gz: $(SRCS)
	$(call build,linux,arm64,)
	$(call tar,linux,arm64)

windows: build/windows_386.zip build/windows_amd64.zip

build/windows_386.zip: $(SRCS)
	$(call build,windows,386,.exe)
	$(call zip,windows,386,.exe)

build/windows_amd64.zip: $(SRCS)
	$(call build,windows,amd64,.exe)
	$(call zip,windows,amd64,.exe)
