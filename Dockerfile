FROM golang:1.22-alpine as builder
WORKDIR /go/app
COPY . .
RUN go build -v -o app main.go
FROM alpine
COPY --from=builder /go/app/ .
CMD ["/app"]