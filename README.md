### Medium Link

https://medium.com/@mertcakmak2/monitoring-the-golang-app-with-prometheus-grafana-new-relic-and-sentry-fce1ca6980b5

### Generate Swagger Docs

```bash
swag init
```

### Run Tests

```bash
go test -v ./...
```

```bash
go test -cover ./...
```


### Run Tests with Coverage
```bash
go test -covermode=count -coverpkg=./... -coverprofile coverage.out -v ./...
```

```bash
go tool cover -html coverage.out
```

### Generate Mock Files

```bash
go generate ./...
```
