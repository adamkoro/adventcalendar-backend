tidy:
	go mod tidy
run:
	SMTP_AUTH=true SMTP_HOST= SMTP_PORT= SMTP_USER= SMTP_PASSWORD= SMTP_FROM= go run main.go
test:
	go test -v ./...
build:
	go build -ldflags="-s -w" -o email-api-consumer