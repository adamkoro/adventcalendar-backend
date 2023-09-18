# adventcalendar-backend

## Build Requirements (dev environment)
- Golang 1.20.3
- Docker 20.10.25
- Docker-compose 1.29.2

## Components
- PostgreSQL
- Redis
- RabbitMQ
- MariaDB

## Environment variables

### Admin-api Init
- `DB_HOST` - PostgreSQL host (default: `localhost`)
- `DB_PORT` - PostgreSQL port (default: `5432`)
- `DB_USER` - PostgreSQL user (default: `adventcalendar`)
- `DB_PASSWORD` - PostgreSQL password (default: `adventcalendar`)
- `DB_NAME` - PostgreSQL database name (default: `adventcalendar`)
- `DB_SSL_MODE` - PostgreSQL ssl mode (default: `disable`)
- `ADMIN_USERNAME` - Admin username (default: `admin`)
- `ADMIN_EMAIL` - Admin email (default: `admin@admin.local`)
- `ADMIN_PASSWORD` - Admin password (default: `admin`)
### Admin-api
**Admin-api Init environment variables are required**
- `SECRET_KEY` - Secret key for JWT token (default: `secret`)
- `PORT` - Api port (default: `8080`)
- `METRICS_PORT` - Api metrics port (default: `8081`)
- `REDIS_HOST` - Redis host (default: `localhost`)
- `REDIS_PORT` - Redis port (default: `6379`)
- `REDIS_PASSWORD` - Redis password (default: `""`)
- `REDIS_DB` - Redis database (default: `0`)

### Email-api Publisher Init
- `DB_HOST` - MariaDB host (default: `localhost`)
- `DB_PORT` - MariaDB port (default: `3306`)
- `DB_USER` - MariaDB user (default: `adventcalendar`)
- `DB_PASSWORD` - MariaDB password (default: `adventcalendar`)
- `DB_NAME` - MariaDB database name (default: `adventcalendar`)
### Email-api Publisher
**Email-api Publisher Init environment variables are required**
- `PORT` - Api port (default: `8080`)
- `METRICS_PORT` - Api metrics port (default: `8081`)
- `RABBITMQ_HOST` - Rabbitmq host (default: `localhost`)
- `RABBITMQ_PORT` - Rabbitmq port (default: `5672`)
- `RABBITMQ_USER` - Rabbitmq user (default: `guest`)
- `RABBITMQ_PASSWORD` - Rabbitmq password (default: `guest`)
- `RABBITMQ_VHOST` - Rabbitmq vhost (default: `/`)
### Email-api Consumer
- `RABBITMQ_HOST` - Rabbitmq host (default: `localhost`)
- `RABBITMQ_PORT` - Rabbitmq port (default: `5672`)
- `RABBITMQ_USER` - Rabbitmq user (default: `guest`)
- `RABBITMQ_PASSWORD` - Rabbitmq password (default: `guest`)
- `RABBITMQ_VHOST` - Rabbitmq vhost (default: `/`)
- `SMTP_AUTH` - SMTP auth (default: `false`)
- `SMTP_HOST` - SMTP host (default: `localhost`)
- `SMTP_PORT` - SMTP port (default: `25`)
- `SMTP_USER` - SMTP user (default: `""`)
- `SMTP_PASSWORD` - SMTP password (default: `""`)
- `SMTP_FROM` - SMTP from (default: `""`)



## How to run
Everything is in Makefile
### Makefile commands
In [admin-api-init](./admin-api-init/) and [admin-api](./admin-api/) directory
- `make run` - Run application
- `make test` - Run tests
- `make tidy` - Update go.mod and go.sum
- `make build` - Build application binary

### Environments
In [root](./) directory
- `make build-images` - Build docker images
#### Dev environment, run services only
- `make compose-up-dev` - Run docker-compose and create dev services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana)
- `make compose-down-dev` - Stop docker-compose and remove dev services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana)
- `make compose-ps-dev` - Show docker-compose dev processes (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana)
#### Stage environment, run services and dockerized application
- `make compose-up-stage` - Run docker-compose and create stage services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana, Adventcalendat-backend-init, Adventcalendar-backend)
- `make compose-down-stage` - Stop docker-compose and remove stage services (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana, Adventcalendat-backend-init, Adventcalendar-backend)
- `make compose-ps-stage` - Show docker-compose stage processes (PostgreSQL, Redis, Rabbitmq, Prometheus, Grafana, Adventcalendat-backend-init, Adventcalendar-backend)

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
## Admin Api endpoints
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

## Email Api endpoints
**TODO**: Swagger documentation
### Public endpoints
- `GET /api/ping`
  - Health check

### Private endpoints - Authentication required
Authentication is required for all endpoints below.
Based on JWT(Json Web Token) authentication. Before using private endpoints, you need to get JWT token from `/api/login` endpoint.

- `GET /api/admin/email` 
  - Get all email patterns
- `POST /api/admin/customemail`
  - Create custom email, which is template for email and automatically send to RabbitMQ
  - Payload (example): `{"emailto": "yourname@gmail.com", "subject": "Test subject", "body": "Test body"}`
- `POST /api/admin/sendemail`
  - Send email to RabbitMQ, which is stored in database
  - Payload (example): `{"name": "customemailpattern"}`
- `POST /api/admin/email`
  - Create email pattern, which is stored in database 
  - Payload (example): `{"name": "name": "customemailpattern", "from": "instace1@localhost", "to": "weblist@localhost", "subject": "Test subject", "body": "Test body"}`
- `PUT /api/admin/email`
  - Update email pattern
- `DELETE /api/admin/email`
  - Delete email
  - Payload (example): `{"name": customname}` 