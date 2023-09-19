package main

import (
	"context"
	"log"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	"github.com/adamkoro/adventcalendar-backend/lib/postgres"
)

func main() {
	// Postgres connection
	var db *postgres.Repository
	conn, err := db.Connect(env.GetDbHost(), env.GetDbUser(), env.GetDbPassword(), env.GetDbName(), env.GetDbPort(), env.GetDbSslMode())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	db = postgres.NewRepository(conn, &ctx)
	log.Println("Postgres connection established.")
	// Migrate
	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Postgres migration completed.")
	// Check if admin user exists
	isAdminExists, err := db.GetUser(env.GetAdminUsername())
	if err != nil {
		log.Println("Admin user does not exist.")
		log.Println("Admin user will be created.")
		// Admin user create
		err = db.CreateUser(env.GetAdminUsername(), env.GetAdminEmail(), env.GetAdminPassword())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Admin user created.")
		return
	}
	log.Println("Admin user exists.")
	// Update admin user if exists and password is not match
	if isAdminExists.Username != "" {
		err = db.CheckUserPassword(env.GetAdminUsername(), env.GetAdminPassword())
		if err != nil {
			log.Println("Admin password is not match.")
			log.Println("Admin user password will be updated.")
			err = db.UpdateUser(env.GetAdminUsername(), env.GetAdminEmail(), env.GetAdminPassword())
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Admin user updated.")
		}
	}
	// Close connection
	defer db.Close()
}
