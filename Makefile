test:
	bash dependencies.sh
	go test -v -race ./...

lint:
	bash dependencies.sh
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run

check: lint test

cover:
	bash dependencies.sh
	go test -race -cover -coverprofile=cover.out ./...
	go tool cover -html=cover.out
	cat cover.out >> coverage.txt