module github.com/adamkoro/adventcalendar-backend/auth-api-init

go 1.20

replace github.com/adamkoro/adventcalendar-backend/lib/postgres => ../lib/postgres

require (
	github.com/adamkoro/adventcalendar-backend/lib/env v0.0.0-00010101000000-000000000000
	github.com/adamkoro/adventcalendar-backend/lib/postgres v0.0.0-00010101000000-000000000000
	github.com/common-nighthawk/go-figure v0.0.0-20210622060536-734e95fb86be
	github.com/rs/zerolog v1.30.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	gorm.io/driver/postgres v1.5.2 // indirect
	gorm.io/gorm v1.25.4 // indirect
)

replace github.com/adamkoro/adventcalendar-backend/lib/env => ../lib/env
