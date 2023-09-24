module github.com/adamkoro/adventcalendar-backend/email-api/publisher-init

go 1.20

replace github.com/adamkoro/adventcalendar-backend/lib/mariadb => ../../lib/mariadb

replace github.com/adamkoro/adventcalendar-backend/lib/env => ../../lib/env

require (
	github.com/adamkoro/adventcalendar-backend/lib/env v0.0.0-00010101000000-000000000000
	github.com/adamkoro/adventcalendar-backend/lib/mariadb v0.0.0-00010101000000-000000000000
	github.com/common-nighthawk/go-figure v0.0.0-20210622060536-734e95fb86be
	github.com/rs/zerolog v1.30.0
)

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
	gorm.io/driver/mysql v1.5.1 // indirect
	gorm.io/gorm v1.25.4 // indirect
)
