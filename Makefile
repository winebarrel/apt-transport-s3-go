GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: all
all: vet test build

.PHONY: build
build:
	go build ./cmd/s3

.PHONY: test
test: vet
	go test -v ./...

.PHONY: vet
vet:
	go vet ./...
