run:
	go run ./cmd/shop-api

test:
	go test ./...

integrate_test:
	go test -race -v ./internal/tests/

.DEFAULT_GOAL=run