.PHONY: install test fmt

install:
	@go install github.com/golang/mock/mockgen@latest
	@go install golang.org/x/tools/cmd/goimports@latest

fmt:
	@goimports -w .
	@go fmt ./...

test: install
	@go test -race -count 1 -v ./...

