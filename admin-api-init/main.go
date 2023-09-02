package main

import (
	"log"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	"github.com/adamkoro/adventcalendar-backend/lib/postgres"
)

func main() {
	db, err := postgres.Connect(env.GetDbHost(), env.GetDbUser(), env.GetDbPassword(), env.GetDbName(), env.GetDbPort(), env.GetDbSslMode())
	if err != nil {
		log.Fatal(err)
	}

	defer postgres.Close(db)

}
