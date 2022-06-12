default: test

clean:
	go clean -testcache ./...

test:
	go test ./...

lint:
	golangci-lint run --disable=structcheck

tidy:
	go mod tidy

run:
	go build -mod vendor -o ./out
	./out --color=always

.PHONY: default clean test lint tidy run
