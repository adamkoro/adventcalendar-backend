SOURCE_DIR = ./src

tidy:
	cd $(SOURCE_DIR) && go mod tidy

run:
	cd $(SOURCE_DIR) && go run main.go

test:
	cd $(SOURCE_DIR) && go test -v ./...

docker-run:
	docker build -t adventcalendar-api:latest -f Dockerfile . && \
	docker run -p 8080:8080 -p 8081:8081 adventcalendar-api:latest

compose-up-dev:
	docker-compose -f compose/dev/docker-compose.yml up -d

compose-down-dev:
	docker-compose -f compose/dev/docker-compose.yml down

compose-ps-dev:
	docker-compose -f compose/dev/docker-compose.yml ps

compose-up-staging:
	docker-compose -f compose/staging/docker-compose.yml up -d

compose-down-staging:
	docker-compose -f compose/staging/docker-compose.yml down

compose-ps-staging:
	docker-compose -f compose/staging/docker-compose.yml ps