COVERAGE = coverage.out

all: build

.PHONY: deps
deps:
	go get github.com/mattn/goveralls
	go get github.com/haya14busa/goverage

.PHONY: build
build:
	go build -v

.PHONY: fmt
fmt:
	gofmt -s -w $$(git ls-files | grep -e '\.go$$' | grep -v -e vendor)
	goimports -w $$(git ls-files | grep -e '\.go$$' | grep -v -e vendor)

.PHONY: test
test:
	DB_NAME=go-active-learning-test go test -v ./...

.PHONY: vet
vet:
	go tool vet --all *.go

.PHONY: test-all
test-all: vet test

.PHONY: cover
cover:
	DB_NAME=go-active-learning-test goverage -v -coverprofile=${COVERAGE} ./...
