### Run Tests

go test -v ./...

go test -cover ./...

### Run Tests with Coverage
go test -covermode=count -coverpkg=./... -coverprofile coverage.out -v ./...

go tool cover -html coverage.out


### Generate Mock Files

go generate ./...
