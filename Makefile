.PHONY: dc run test lint

dc:
	docker-compose up  --remove-orphans --build

build:
	go build -race -o app cmd/main.go

run:
	go build -race -o app cmd/main.go && \
	KRAKEN_HOST="https://api.kraken.com" \
	HTTP_ADDR=:8090 \
	DEBUG_ERRORS=1 \
	./app

test:
	go test -race ./...

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4

lint:
	golangci-lint run ./...

generate:
	go generate ./...