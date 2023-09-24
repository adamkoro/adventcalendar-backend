tidy:
	go mod tidy
run:
	LOG_LEVEL=debug DB_HOST=localhost DB_PORT=3306 DB_USER=adventcalendar DB_PASSWORD=adventcalendar DB_NAME=adventcalendar go run main.go
test:
	go test -v ./...
build:
	go build -ldflags="-s -w" -o email-api-publisher-init