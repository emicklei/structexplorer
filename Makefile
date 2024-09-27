lint:
	golangci-lint run

test:
	go test -v -cover -race