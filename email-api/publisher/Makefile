tidy:
	go mod tidy
run:
	PORT=8082 METRICS_PORT=8083 go run main.go
test:
	go test -v ./...
build:
	go build -ldflags="-s -w" -o email-api-publisher