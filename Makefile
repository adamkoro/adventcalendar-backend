API_IMAGE_NAME="adventcalendar-backend"
API_IMAGE_TAG="latest"
INIT_IMAGE_NAME="adventcalendar-backend-init"
INIT_IMAGE_TAG="latest"

compose-up-dev:
	docker-compose -f compose/dev/docker-compose.yml up -d
compose-down-dev:
	docker-compose -f compose/dev/docker-compose.yml down
compose-ps-dev:
	docker-compose -f compose/dev/docker-compose.yml ps
compose-up-staging:
	docker-compose -f compose/staging/docker-compose.yml up -d --build
compose-down-staging:
	docker-compose -f compose/staging/docker-compose.yml down
compose-ps-staging:
	docker-compose -f compose/staging/docker-compose.yml ps

build-images:
	docker build -t $(API_IMAGE_NAME):${API_IMAGE_TAG} -f admin-api/Dockerfile .
	docker build -t $(INIT_IMAGE_NAME):${INIT_IMAGE_TAG} -f admin-api-init/Dockerfile .