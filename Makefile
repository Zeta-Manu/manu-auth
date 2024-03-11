run:
	go run cmd/main.go

tidy:
	go mod tidy

vendor:
	go mod vendor

swag:
	swag init -g cmd/main.go

docker-build:
	docker build --tag manu-auth .
