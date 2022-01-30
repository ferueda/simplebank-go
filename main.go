package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ferueda/simplebank-go/api"
	db "github.com/ferueda/simplebank-go/db/sqlc"
	"github.com/ferueda/simplebank-go/token"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var dbAddr string
var appAddr string
var dbDriver string
var tokenKey string

func init() {
	env := os.Getenv("ENV")
	if env == "dev" || env == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	dbAddr = os.Getenv("DB_HOST")
	appAddr = os.Getenv("APP_HOST")
	dbDriver = os.Getenv("DB_DRIVER")
	tokenKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
}

func main() {
	conn, err := sql.Open(dbDriver, dbAddr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	tm, err := token.NewPasetoMaker(tokenKey)
	if err != nil {
		log.Fatal("cannot create token maker: %w", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, tm)
	if err != nil {
		log.Fatal("cannot create server: %w", err)
	}

	if err = server.Start(appAddr); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
