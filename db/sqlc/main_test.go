package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/backendmaster/simple_bank/util"
	_ "github.com/lib/pq"
)

var testQuires *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatalf("Load Config Failed: %v", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't not connect to database")
	}
	testQuires = New(testDB)

	os.Exit(m.Run())
}
