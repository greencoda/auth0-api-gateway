.PHONY: all build clean deps docker-build mock run test
.SILENT: swag test test-100 test-cover

BINARY_NAME=auth0-api-gateway
TEST_PACKAGES=$(shell go list ./... | grep -v '/mocks/\|/tools/\|/cmd\|/module')

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
	go test ${TEST_PACKAGES} -cover -count=1

test-100:
	go test ${TEST_PACKAGES} -cover -count=100
	
test-cover:
	go test ${TEST_PACKAGES} -cover -coverprofile=coverage.out && go tool cover -html=coverage.out