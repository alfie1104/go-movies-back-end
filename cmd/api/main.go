package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
	DSN string
	Domain string
	DB *sql.DB
}

func main() {
	// set application config
	var app application

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " sslmode=disable timezone=UTC connect_timeout=5","Postgres connection string")
	flag.Parse()

	// connect to the database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = conn

	app.Domain = "example.com"

	log.Println("Starting application on port", port)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())

	if err != nil {
		log.Fatal(err)
	}
}