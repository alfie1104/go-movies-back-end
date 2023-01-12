package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
	DSN string
	Domain string
	DB repository.DatabaseRepo
	auth Auth
	JWTSecret string
	JWTIssuer string
	JWTAudience string
	CookieDomain string
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
	flag.StringVar(&app.JWTSecret, "jwt-secret","verysecret","signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer","example.com","signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience","example.com","signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain","localhost","cookie domain")
	flag.StringVar(&app.Domain, "domain","example.com","domain")
	flag.Parse()

	// connect to the database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.auth = Auth{
		Issuer: app.JWTIssuer,
		Audience: app.JWTAudience,
		Secret: app.JWTSecret,
		TokenExpiry: time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath: "/",
		CookieName: "__Host-refresh_token",
		CookieDomain: app.CookieDomain,
	}

	log.Println("Starting application on port", port)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())

	if err != nil {
		log.Fatal(err)
	}
}