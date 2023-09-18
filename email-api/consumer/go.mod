module github.com/adamkoro/adventcalendar-backend/email-api/consumer

go 1.20

replace github.com/adamkoro/adventcalendar-backend/lib/rabbitmq => ../../lib/rabbitmq

require (
	github.com/adamkoro/adventcalendar-backend/lib/env v0.0.0-00010101000000-000000000000
	github.com/adamkoro/adventcalendar-backend/lib/model v0.0.0-00010101000000-000000000000
	github.com/adamkoro/adventcalendar-backend/lib/rabbitmq v0.0.0-00010101000000-000000000000
	github.com/rabbitmq/amqp091-go v1.8.1
)

replace github.com/adamkoro/adventcalendar-backend/lib/env => ../../lib/env

replace github.com/adamkoro/adventcalendar-backend/lib/model => ../../lib/model
