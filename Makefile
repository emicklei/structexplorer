lint:
	golangci-lint run

test: lint
	go test -v -cover -race