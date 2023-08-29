# adventcalendar-backend

## Build Requirements (dev environment)
- Golang 1.20.3
- Docker 20.10.25
- Docker-compose 1.29.2

## Components
- PostgreSQL
- Redis
- Rabbitmq

## How to run
Everything is in Makefile
### Makefile commands
- `make run` - Run application
- `make test` - Run tests
- `make tidy` - Update go.mod and go.sum
- `make docker-run` - Run application in docker
- `make compose-up` - Run docker-compose and create dev services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana)
- `make compose-down` - Stop docker-compose and remove dev services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana)
- `make compose-ps` - Show docker-compose dev processes (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana)

#### Before run
```shell
make compose-up
```
#### Shell
```shell
make run
```
#### Docker
```shell
make docker-run
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
