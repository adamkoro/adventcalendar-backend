package main

import (
	"context"
	"log"
	"time"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	md "github.com/adamkoro/adventcalendar-backend/lib/mariadb"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	figure.NewFigure("AdventCalendar Email Publisher", "big", false).Print()
	// MariaDB connection
	log.Println("Connecting to the MariaDB...")
	var db *md.Repository
	conn, err := db.Connect(env.GetDbUser(), env.GetDbPassword(), env.GetDbHost(), env.GetDbName(), env.GetDbPort())
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db = md.NewRepository(conn, &ctx)
	log.Println("Connected to the MariaDB.")
	// MariaDB migration
	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	// Check if default email exists
	_, err = db.GetEmailByName("default")
	if err != nil {
		log.Println("Failed to get default email.")
		log.Println("Default email will be created.")
		var email md.Email
		email.Name = "default"
		email.From = "adventcalendar@localhost"
		email.To = "adventcalendar@localhost"
		email.Subject = "Advent Calendar"
		email.Body = "Hello, World!"
		err = db.CreateEmail(&email)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Default email has been created.")
		return
	}
	log.Println("Default email exists.")
	// Close MariaDB connection
	defer db.Close()
}
