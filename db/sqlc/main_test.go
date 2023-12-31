package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

var cleanup func()

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:postgres123@localhost:5432/simple_bank_test?sslmode=disable"
)

func TestMain(m *testing.M){
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	
	cleanup = func() {
		testQueries.db.ExecContext(context.Background(), "TRUNCATE TABLE transfers")
		testQueries.db.ExecContext(context.Background(), "TRUNCATE TABLE entries")
		_, err2 := testQueries.db.ExecContext(context.Background(), "TRUNCATE TABLE accounts CASCADE")
		if err2 != nil {
			log.Fatal("cannot truncate accounts: ", err2)
		}
	}
	os.Exit(m.Run())
}
