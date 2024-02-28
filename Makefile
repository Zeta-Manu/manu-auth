run:
	go run cmd/main.go

tidy:
	go mod tidy

vendor:
	go mod vendor

docker:build:
	docker build --tag manu-auth .
