.PHONY: dc run test lint

dc:
	docker-compose up  --remove-orphans --build

run:
	go build -o app cmd/url-shortener/main.go && CONFIG_PATH=internal/config/config.yaml ./app

test:
	go test -race ./...

lint:
	golangci-lint run

