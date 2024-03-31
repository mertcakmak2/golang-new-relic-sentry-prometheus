### Run Tests

go test -v ./...

go test -cover ./...

### Run Tests with Coverage
go test -covermode=count -coverpkg=./... -coverprofile coverage.out -v ./...

go tool cover -html coverage.out


### Generate Mock Files

go generate ./...

### Medium Link

https://medium.com/@mertcakmak2/monitoring-the-golang-app-with-prometheus-grafana-new-relic-and-sentry-fce1ca6980b5