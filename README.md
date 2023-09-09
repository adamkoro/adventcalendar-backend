# adventcalendar-backend

## Build Requirements (dev environment)
- Golang 1.20.3
- Docker 20.10.25
- Docker-compose 1.29.2

## Components
- PostgreSQL
- Redis
- Rabbitmq

## Environment variables
- `SECRET_KEY` - Secret key for JWT token (default: `secret`)
- `ADMIN_USERNAME` - Admin username (default: `admin`)
- `ADMIN_EMAIL` - Admin email (default: `admin@admin.local`)
- `ADMIN_PASSWORD` - Admin password (default: `admin`)
- `PORT` - Api port (default: `8080`)
- `METRICS_PORT` - Api metrics port (default: `8081`)
- `DB_HOST` - PostgreSQL host (default: `localhost`)
- `DB_PORT` - PostgreSQL port (default: `5432`)
- `DB_USER` - PostgreSQL user (default: `postgres`)
- `DB_PASSWORD` - PostgreSQL password (default: `postgres`)
- `DB_NAME` - PostgreSQL database name (default: `postgres`)
- `DB_SSL_MODE` - PostgreSQL ssl mode (default: `disable`)
- `REDIS_HOST` - Redis host (default: `localhost`)
- `REDIS_PORT` - Redis port (default: `6379`)
- `REDIS_PASSWORD` - Redis password (default: `""`)
- `REDIS_DB` - Redis database (default: `0`)
- `RABBITMQ_HOST` - Rabbitmq host (default: `localhost`)
- `RABBITMQ_PORT` - Rabbitmq port (default: `5672`)
- `RABBITMQ_USER` - Rabbitmq user (default: `guest`)
- `RABBITMQ_PASSWORD` - Rabbitmq password (default: `guest`)
- `RABBITMQ_VHOST` - Rabbitmq vhost (default: `/`)


## How to run
Everything is in Makefile
### Makefile commands
- `make run` - Run application
- `make test` - Run tests
- `make tidy` - Update go.mod and go.sum
- `make build` - Build application binary

#### Dev environment run services only
- `make compose-up-dev` - Run docker-compose and create dev services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana)
- `make compose-down-dev` - Stop docker-compose and remove dev services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana)
- `make compose-ps-dev` - Show docker-compose dev processes (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana)
#### Stage environment run services and dockerized application
- `make compose-up-stage` - Run docker-compose and create stage services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana, Adventcalendar-backend)
- `make compose-down-stage` - Stop docker-compose and remove stage services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana, Adventcalendar-backend)
- `make compose-ps-stage` - Show docker-compose stage processes (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana, Adventcalendar-backend)

### Create dev environment
```shell
make compose-up-dev
```
### Setup application
Before run application, you need to initialize database and create admin user.
#### Initialize database
```shell
cd admin-api-init && make run
```
#### Run application
```shell
cd admin-api && make run
```
### Access to application
- `http://localhost:8080` - Application
- `http://localhost:8081` - Prometheus metrics
## Api endpoints
**TODO**: Swagger documentation
### Public endpoints
- `GET /api/ping`
  - Health check
- `POST /api/login` 
  - User login
  - Payload (example): `{"username": "testuser1", "password": "testpassword1"}`
- `POST /api/logout` 
  - User logout
- `GET /metrics`
  - Prometheus metrics
### Private endpoints - Authentication required
Authentication is required for all endpoints below.
Based on JWT(Json Web Token) authentication. Before using private endpoints, you need to get JWT token from `/api/login` endpoint.

- `GET /api/admin/user` 
  - Get user
  - Payload (example): `{"username": "testuser1"}`
- `POST /api/admin/user` 
  - Create user
  - Payload (example): `{"username": "testuser1", "email": "testuser1@test.local", "password": "testpassword1"}`
- `PUT /api/admin/user` 
  - Update user
  - Payload (example): `{"username": "testuser1", "email": "testuser1@gmail.com", "password": "testpassword1"}`
- `DELETE /api/admin/user`
  - Delete user
  - Payload (example): `{"username": "testuser1"}`
- `GET /api/admin/users` 
  - Get all users
