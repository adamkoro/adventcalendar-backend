# adventcalendar-backend
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![MariaDB](https://img.shields.io/badge/MariaDB-003545?style=for-the-badge&logo=mariadb&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)

![Suse](https://img.shields.io/badge/SUSE-0C322C?style=for-the-badge&logo=SUSE&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)

[![Build Status](https://drone.adamkoro.com/api/badges/adamkoro/adventcalendar-backend/status.svg?ref=refs/heads/main)](https://drone.adamkoro.com/adamkoro/adventcalendar-backend)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)]()

| Component | Code report |
| ------ | ------ |
| Admin Api | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/admin-api)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/admin-api) |
| Auth Api Init | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/auth-api-init)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/auth-api-init) |
| Auth Api | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/auth-api)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/auth-api) |
| Publisher Init | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/email-api/publisher-init)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/email-api/publisher-init) |
| Publisher | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/email-api/publisher)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/email-api/publisher) |
| Consumer | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/email-api/consumer)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/email-api/consumer) |
| Public Api | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/public-api)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/email-api/public-api) |

| Library | Code report |
| ------ | ------ |
| Env | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/lib/env)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/lib/env) |
| JWT | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/lib/jwt)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/lib/jwt) |
| MariaDB | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/lib/mariadb)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/lib/mariadb) |
| MongoDB | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/lib/mongo)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/lib/mongo) |
| Postgres | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/lib/postgres)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/lib/postgres) |
| RabbitMQ | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/lib/rabbitmq)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/lib/rabbitmq) |
| Redis | [![Go Report Card](https://goreportcard.com/badge/github.com/adamkoro/adventcalendar-backend/lib/redis)](https://goreportcard.com/report/github.com/adamkoro/adventcalendar-backend/lib/redis) |

## Build Requirements (dev environment)
- Golang 1.20.3
- Docker 20.10.25
- Docker-compose 1.29.2

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

## Auth Api endpoints
**TODO**: Swagger documentation
### Public endpoints
- `GET /api/ping`
  - Health check
- `GET /metrics`
  - Prometheus metrics
- `POST /api/auth/login` 
  - User login
  - Payload (example): `{"username": "testuser1", "password": "testpassword1"}`
- `POST /api/auth/logout` 
  - User logout
## Admin Api endpoints
**TODO**: Swagger documentation
### Public endpoints
- `GET /api/ping`
  - Health check
- `GET /metrics`
  - Prometheus metrics
### Private endpoints - Authentication required
Authentication is required for all endpoints below.
Based on JWT(Json Web Token) authentication. Before using private endpoints, you need to get JWT token from `/api/auth/login` endpoint.

- `GET /api/admin/usermanage/user` 
  - Get user
  - Payload (example): `{"username": "testuser1"}`
- `POST /api/admin/usermanage/user` 
  - Create user
  - Payload (example): `{"username": "testuser1", "email": "testuser1@test.local", "password": "testpassword1"}`
- `PUT /api/admin/usermanage/user` 
  - Update user
  - Payload (example): `{"username": "testuser1", "email": "testuser1@gmail.com", "password": "testpassword1"}`
- `DELETE /api/admin/usermanage/user`
  - Delete user
  - Payload (example): `{"username": "testuser1"}`
- `GET /api/admin/usermanage/users` 
  - Get all users

## Email Api endpoints
**TODO**: Swagger documentation
### Public endpoints
- `GET /api/ping`
  - Health check
- `GET /metrics`
  - Prometheus metrics

### Private endpoints - Authentication required
Authentication is required for all endpoints below.
Based on JWT(Json Web Token) authentication. Before using private endpoints, you need to get JWT token from `/api/auth/login` endpoint.

- `GET /api/admin/emailmanage/email` 
  - Get all email patterns
- `POST /api/admin/emailmanage/customemail`
  - Create custom email, which is template for email and automatically send to RabbitMQ
  - Payload (example): `{"emailto": "yourname@gmail.com", "subject": "Test subject", "body": "Test body"}`
- `POST /api/admin/emailmanage/sendemail`
  - Send email to RabbitMQ, which is stored in database
  - Payload (example): `{"name": "customemailpattern"}`
- `POST /api/admin/emailmanage/email`
  - Create email pattern, which is stored in database 
  - Payload (example): `{"name": "name": "customemailpattern", "from": "instace1@localhost", "to": "weblist@localhost", "subject": "Test subject", "body": "Test body"}`
- `PUT /api/admin/emailmanage/email`
  - Update email pattern
- `DELETE /api/admin/emailmanage/email`
  - Delete email
  - Payload (example): `{"name": customname}` 