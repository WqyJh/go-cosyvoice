
.PHONY: test
test:
	go test -race -v ./...

.PHONY: cov
cov:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run
