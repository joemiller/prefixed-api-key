deps:
	@go get

lint:
	@golangci-lint run -v --timeout=3m

test:
	@go test -v -race -coverprofile=cover.out ./...

cov:
	@go tool cover -html=cover.out

.PHONY: deps lint test cov