API_SOURCE_DIR = ./admin-api

tidy-admin-api:
	cd $(API_SOURCE_DIR) && go mod tidy

run-admin-api:
	cd $(API_SOURCE_DIR) && go run main.go

test-admin-api:
	cd $(API_SOURCE_DIR) && go test -v ./...


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