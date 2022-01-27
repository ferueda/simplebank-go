package db

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

var dbAddr string
var dbDriver string

func init() {
	env := os.Getenv("ENV")
	if env == "dev" || env == "" {
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}

	dbAddr = os.Getenv("DB_HOST")
	dbDriver = os.Getenv("DB_DRIVER")

	rand.Seed(time.Now().UnixNano())
}

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbAddr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
