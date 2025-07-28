.PHONY: all build clean deps lint mock run test test-cover
.SILENT: run test test-cover

BINARY_NAME=auth0-api-gateway
TESTABLE_PACKAGES=$(shell go list ./... | grep -v '/mocks/\|/cmd\|/module')

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/$(BINARY_NAME) cmd/main.go 

clean:
	rm -f bin/*
	rm -f *.zip
	rm -f *~
	find . -name "*.coverprofile" | xargs rm
	find . -name "coverage.out" | xargs rm
	go clean -cache -testcache -i

deps:
	go mod tidy
	go mod vendor
	go mod download -x

lint: deps
	golangci-lint run -v

mock:
	rm -rf internal/mocks/*
	mockery --all --dir=internal --output=internal/mocks --keeptree

run:
	go run cmd/main.go

test:
	go test ${TESTABLE_PACKAGES} -cover -count=1 -cover -coverprofile=coverage.out
	
test-cover:
	go test ${TESTABLE_PACKAGES} -cover -coverprofile=coverage.out && go tool cover -html=coverage.out