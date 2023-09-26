package main

import (
	"context"
	"time"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	md "github.com/adamkoro/adventcalendar-backend/lib/mariadb"
	"github.com/common-nighthawk/go-figure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	db          md.Repository
	isConnected bool
	email       md.Email
)

func main() {
	figure.NewFigure("AdventCalendar Email Publisher Init", "big", false).Print()
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
	// Mariadb connection check
	/////////////////////////
	for {
		log.Info().Msg("establishing connection to the mariadb...")
		mariadbConn, err := db.Connect(env.GetDbUser(), env.GetDbPassword(), env.GetDbHost(), env.GetDbName(), env.GetDbPort())
		if err != nil {
			log.Debug().Msg(err.Error())
		}
		ctx := context.Background()
		db := md.NewRepository(mariadbConn, &ctx)
		log.Debug().Msg("pinging the mariadb...")
		err = db.Ping()
		if err != nil {
			log.Info().Msg("lost connection to the mariadb, reconnecting...")
			log.Error().Msg(err.Error())
			isConnected = false
		} else {
			isConnected = true
			log.Debug().Msg("pinging the mariadb successful")
			log.Info().Msg("establishing connection to the mariadb successful")
			if !isConnected {
				log.Info().Msg("reconnected to the mariadb")
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
			// Check if default email exists if not create it
			/////////////////////////
			log.Debug().Msg("checking if default email exists...")
			_, err = db.GetEmailByName(&md.EmailRequest{Name: "default"})
			if err != nil {
				log.Info().Msg("default email does not exist")
				log.Info().Msg("creating default email...")
				email.Name = "default"
				email.From = "adventcalendar@localhost"
				email.To = "adventcalendar@localhost"
				email.Subject = "Advent Calendar"
				email.Body = "Hello, World!"
				log.Debug().Msg("creating default email in database...")
				err = db.CreateEmail(&email)
				if err != nil {
					log.Fatal().Msg(err.Error())
				}
				log.Info().Msg("default email created")
				return
			} else {
				log.Info().Msg("default email exists")
			}
			log.Info().Msg("all required actions completed successfully")
			log.Debug().Msg("closing connection to the mariadb...")
			/////////////////////////
			// Close connection
			/////////////////////////
			db.Close()
			log.Debug().Msg("closing connection to the mariadb successful")
			break
		}
	}
}
