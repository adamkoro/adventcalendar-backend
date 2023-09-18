tidy:
	go mod tidy
run:
	PORT=8082 METRICS_PORT=8083 DB_HOST=localhost DB_PORT=5432 DB_USER=adventcalendar DB_PASSWORD=adventcalendar DB_NAME=adventcalendar go run main.go
test:
	go test -v ./...
build:
	go build -ldflags="-s -w" -o email-api-publisher