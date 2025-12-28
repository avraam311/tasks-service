.PHONY: lint, up, test

lint:
	go vet ./...
	golangci-lint run ./...

up:
	docker build -t tasks-service .
	docker run -p 8080:8080 tasks-service

test:
	go test ./...