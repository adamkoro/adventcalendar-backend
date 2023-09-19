tidy:
	go mod tidy
run:
	PORT=8088 METRICS_PORT=8089 DB_HOST=localhost DB_PORT=27017 DB_USER=adventcalendar DB_PASSWORD=adventcalendar DB_NAME=adventcalendar go run main.go
test:
	go test -v ./...
build:
	go build -ldflags="-s -w" -o public-api