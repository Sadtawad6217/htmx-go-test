package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func Connect() *sqlx.DB {
	viper.SetConfigFile("cmd/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	dbURL := viper.GetString("database.url")
	if dbURL == "" {
		log.Fatalf("Database URL is not set in config.yaml")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %s", err)
	}

	fmt.Println("Successfully connected to database!")
	return db
}
