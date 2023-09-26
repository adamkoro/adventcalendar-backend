package main

import (
	"context"
	"time"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/common-nighthawk/go-figure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	db          pg.Repository
	isConnected bool
)

func main() {
	figure.NewFigure("AdventCalendar Auth Api Init", "big", false).Print()
	/////////////////////////
	// Logger setup
	/////////////////////////
	logLevel := env.GetLogLevel()
	switch logLevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	/////////////////////////
	// Postgres connection check
	/////////////////////////
	for {
		log.Info().Msg("establishing connection to the postgres...")
		postgresConn, err := db.Connect(env.GetDbHost(), env.GetDbUser(), env.GetDbPassword(), env.GetDbName(), env.GetDbPort(), env.GetDbSslMode())
		if err != nil {
			log.Error().Msg(err.Error())
		}
		ctx := context.Background()
		db := pg.NewRepository(postgresConn, &ctx)
		log.Debug().Msg("pinging the postgres...")
		err = db.Ping()
		if err != nil {
			log.Info().Msg("lost connection to the postgres, reconnecting...")
			log.Error().Msg(err.Error())
			isConnected = false
		} else {
			isConnected = true
			log.Debug().Msg("pinging the postgres successful")
			log.Info().Msg("establishing connection to the postgres successful")
			if !isConnected {
				log.Info().Msg("reconnected to the postgres")
				isConnected = true
			}
		}
		time.Sleep(5 * time.Second)
		if isConnected {
			/////////////////////////
			// Database migration
			/////////////////////////
			log.Debug().Msg("database struct migration started...")
			err := db.Migrate()
			if err != nil {
				log.Fatal().Msg(err.Error())
			}
			log.Info().Msg("database struct migration successful")
			/////////////////////////
			// Check if admin user exists
			/////////////////////////
			log.Debug().Msg("checking admin user exists...")
			isAdminExists, err := db.GetUser(&pg.UserRequest{Username: env.GetAdminUsername()})
			if err != nil {
				log.Info().Msg("admin user does not exists")
				log.Info().Msg("admin user creating...")
				// Admin user create
				err = db.CreateUser(&pg.CreateUserRequest{Username: env.GetAdminUsername(), Email: env.GetAdminEmail(), Password: env.GetAdminPassword()})
				if err != nil {
					log.Fatal().Msg(err.Error())
				}
				log.Info().Msg("admin user created")
				return
			}
			log.Info().Msg("admin user is exists")
			/////////////////////////
			// Update admin user if exists and password is not match
			/////////////////////////
			if isAdminExists.Username != "" {
				log.Debug().Msg("checking admin user password...")
				err = db.CheckUserPassword(&pg.LoginRequest{Username: env.GetAdminUsername(), Password: env.GetAdminPassword()})
				if err != nil {
					log.Info().Msg("admin user password is not match")
					log.Info().Msg("admin user password updating...")
					err = db.UpdateUser(&pg.UpdateUserRequest{Username: env.GetAdminUsername(), Email: env.GetAdminEmail(), Password: env.GetAdminPassword()})
					if err != nil {
						log.Fatal().Msg(err.Error())
					}
					log.Info().Msg("admin user password updated")
				}
				log.Info().Msg("admin user password is match")
				log.Debug().Msg("checking admin user password successful")
			}
			log.Info().Msg("all required actions completed successfully")
			log.Debug().Msg("closing connection to the postgres...")
			/////////////////////////
			// Close connection
			/////////////////////////
			db.Close()
			log.Debug().Msg("closing connection to the postgres successful")
			break
		}
	}
}
