## Install dependencies
setup:
	go mod download

build:
	go fmt ./...
	GOOS=$(PLATFORM) GOARCH=$(GOARCH) go build -o ezinsights cmd/ezinsights/main.go

## Analysis exec vet and lint
analysis: setup
	go vet ./...
	go get -u golang.org/x/lint/golint
	golint -set_exit_status $$(go list ./... | grep -v /vendor/)

## Run tests
test: setup
	go test -race ./...

.PHONY: setup build analysis test
