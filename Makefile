run:
	go run ./cmd/shop-api

test:
	go test ./...

integrate_test:
	curl -X GET localhost:8080/items
	curl -X GET localhost:8080/items/1
	curl -X DELETE localhost:8080/items/1
	curl -X POST --data '{"name": "Item 1", "price": 0}' localhost:8080/items
	curl -X POST --data '{"name": "Item 2"}' localhost:8080/items
	curl -X GET localhost:8080/items
	curl -X GET localhost:8080/items/1
	curl -X DELETE localhost:8080/items/1
	curl -X GET localhost:8080/items

.DEFAULT_GOAL=run