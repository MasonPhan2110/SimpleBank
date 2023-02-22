package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const (
	DB_DRIVER = "postgres"
	DB_SOURCE = "postgresql://root:secret@localhost:3000/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	// config, err := util.LoadConfig("../..")
	// if err != nil {
	// 	log.Fatal("cannot cload config: ", err)
	// }
	testDB, err = sql.Open(DB_DRIVER, DB_DRIVER)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
