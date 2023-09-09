package main

import (
	"log"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	"github.com/adamkoro/adventcalendar-backend/lib/postgres"
)

func main() {
	// Postgres connection
	db, err := postgres.Connect(env.GetDbHost(), env.GetDbUser(), env.GetDbPassword(), env.GetDbName(), env.GetDbPort(), env.GetDbSslMode())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Postgres connection established.")
	// Migrate
	err = postgres.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Postgres migration completed.")
	// Admin user create
	err = postgres.CreateUser(db, env.GetAdminUsername(), env.GetAdminEmail(), env.GetAdminPassword())
	if err != nil {
		log.Println(err)
	}
	log.Println("Admin user created.")
	// Check if admin user exists
	isAdminExists, err := postgres.GetUser(db, env.GetAdminUsername())
	if err != nil {
		log.Println(err)
	}
	log.Println("Admin user exists.")
	// Update admin user if exists
	if isAdminExists.Username != "" {
		err = postgres.UpdateUser(db, env.GetAdminUsername(), env.GetAdminEmail(), env.GetAdminPassword())
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("Admin user updated.")
	// Close connection
	defer postgres.Close(db)
}