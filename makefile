# generates mock services
mock:
	go generate ./...

# generates swagger api docs
api-docs:
	swag init

# runs tests
run-tests:
	go clean -cache
	go test -v ./...

build:
	go build -o bin/main main.go

docker-build:
	docker build -t mertcakmak2/go-e2e .

run:
	go run main.go
