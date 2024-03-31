hello:
	echo "Hello"

api-docs:
	swag init

unit-test:
	go clean -cache
	go test -v ./...

build:
	go build -o bin/main main.go

docker-build:
	docker build -t mertcakmak2/go-e2e .

run:
	go run main.go
