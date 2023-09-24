package env

import (
	"os"
	"strconv"
)

func GetLogLevel() string {
	logLevel := "info"
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	return logLevel
}

func GetSecretKey() string {
	secret := "secret"
	if os.Getenv("SECRET_KEY") != "" {
		secret = os.Getenv("SECRET_KEY")
	}
	return secret
}

func GetAdminUsername() string {
	username := "admin"
	if os.Getenv("ADMIN_USERNAME") != "" {
		username = os.Getenv("ADMIN_USERNAME")
	}
	return username
}

func GetAdminEmail() string {
	email := "admin@admin.local"
	if os.Getenv("ADMIN_EMAIL") != "" {
		email = os.Getenv("ADMIN_EMAIL")
	}
	return email
}

func GetAdminPassword() string {
	password := "admin"
	if os.Getenv("ADMIN_PASSWORD") != "" {
		password = os.Getenv("ADMIN_PASSWORD")
	}
	return password
}

func GetHttpPort() int {
	port := 8080
	if os.Getenv("PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}
	return port
}

func GetMetricsPort() int {
	port := 8081
	if os.Getenv("METRICS_PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("METRICS_PORT"))
	}
	return port
}

func GetDbHost() string {
	host := "localhost"
	if os.Getenv("DB_HOST") != "" {
		host = os.Getenv("DB_HOST")
	}
	return host
}

func GetDbPort() int {
	port := 5432
	if os.Getenv("DB_PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	}
	return port
}

func GetDbUser() string {
	user := "adventcalendar"
	if os.Getenv("DB_USER") != "" {
		user = os.Getenv("DB_USER")
	}
	return user
}

func GetDbPassword() string {
	password := "adventcalendar"
	if os.Getenv("DB_PASSWORD") != "" {
		password = os.Getenv("DB_PASSWORD")
	}
	return password
}

func GetDbName() string {
	name := "adventcalendar"
	if os.Getenv("DB_NAME") != "" {
		name = os.Getenv("DB_NAME")
	}
	return name
}

func GetDbSslMode() string {
	sslMode := "disable"
	if os.Getenv("DB_SSL_MODE") != "" {
		sslMode = os.Getenv("DB_SSL_MODE")
	}
	return sslMode
}

func GetRedisHost() string {
	host := "localhost"
	if os.Getenv("REDIS_HOST") != "" {
		host = os.Getenv("REDIS_HOST")
	}
	return host
}

func GetRedisPort() int {
	port := 6379
	if os.Getenv("REDIS_PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("REDIS_PORT"))
	}
	return port
}

func GetRedisPassword() string {
	password := ""
	if os.Getenv("REDIS_PASSWORD") != "" {
		password = os.Getenv("REDIS_PASSWORD")
	}
	return password
}

func GetRedisDb() int {
	db := 0
	if os.Getenv("REDIS_DB") != "" {
		db, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
	}
	return db
}

func GetRabbitmqHost() string {
	host := "localhost"
	if os.Getenv("RABBITMQ_HOST") != "" {
		host = os.Getenv("RABBITMQ_HOST")
	}
	return host
}

func GetRabbitmqPort() int {
	port := 5672
	if os.Getenv("RABBITMQ_PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("RABBITMQ_PORT"))
	}
	return port
}

func GetRabbitmqUser() string {
	user := "rabbitmq"
	if os.Getenv("RABBITMQ_USER") != "" {
		user = os.Getenv("RABBITMQ_USER")
	}
	return user
}

func GetRabbitmqPassword() string {
	password := "rabbitmq"
	if os.Getenv("RABBITMQ_PASSWORD") != "" {
		password = os.Getenv("RABBITMQ_PASSWORD")
	}
	return password
}

func GetRabbitmqVhost() string {
	vhost := "/"
	if os.Getenv("RABBITMQ_VHOST") != "" {
		vhost = os.Getenv("RABBITMQ_VHOST")
	}
	return vhost
}

func GetSmtpAuth() bool {
	auth := false
	if os.Getenv("SMTP_AUTH") != "" {
		auth, _ = strconv.ParseBool(os.Getenv("SMTP_AUTH"))
	}
	return auth
}

func GetSmtpHost() string {
	host := "localhost"
	if os.Getenv("SMTP_HOST") != "" {
		host = os.Getenv("SMTP_HOST")
	}
	return host
}

func GetSmtpPort() string {
	port := "25"
	if os.Getenv("SMTP_PORT") != "" {
		port = os.Getenv("SMTP_PORT")
	}
	return port
}

func GetSmtpUser() string {
	user := ""
	if os.Getenv("SMTP_USER") != "" {
		user = os.Getenv("SMTP_USER")
	}
	return user
}

func GetSmtpPassword() string {
	password := ""
	if os.Getenv("SMTP_PASSWORD") != "" {
		password = os.Getenv("SMTP_PASSWORD")
	}
	return password
}

func GetSmtpFrom() string {
	from := "adventcalendar@localhost"
	if os.Getenv("SMTP_FROM") != "" {
		from = os.Getenv("SMTP_FROM")
	}
	return from
}
