## Install dependencies
setup:
	dep version > /dev/null || curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep check || dep ensure

## Analysis exec vet and lint
analysis: setup
	go vet ./...
	go get -u golang.org/x/lint/golint
	golint -set_exit_status $$(go list ./... | grep -v /vendor/)

## Run tests
test: setup
	go test -race ./...

.PHONY: setup analysis test
