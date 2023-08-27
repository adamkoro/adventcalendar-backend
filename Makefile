SOURCE_DIR = ./src

tidy:
	cd $(SOURCE_DIR) && go mod tidy

run:
	cd $(SOURCE_DIR) && go run main.go

test:
	cd $(SOURCE_DIR) && go test -v ./...

docker-run:
	docker build -t adventcalendar-api:latest -f Dockerfile . && \
	docker run -p 8080:8080 adventcalendar-api:latest