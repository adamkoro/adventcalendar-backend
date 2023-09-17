module github.com/adamkoro/adventcalendar-backend/lib/rabbitmq

go 1.20

require (
	github.com/adamkoro/adventcalendar-backend/lib/model v0.0.0-00010101000000-000000000000
	github.com/rabbitmq/amqp091-go v1.8.1
)

replace github.com/adamkoro/adventcalendar-backend/lib/model => ../model
