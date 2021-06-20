.DEFAULT_GOAL := compile

compile: lint
	@go build -o crawler cmd/main.go

lint:
	 golangci-lint run ./...