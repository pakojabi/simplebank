package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pakojabi/simplebank/util"
)

var testQueries *Queries
var testDB *sql.DB

var cleanup func()

func TestMain(m *testing.M){
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.TestDBSource)
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
		_, err3 := testQueries.db.ExecContext(context.Background(), "TRUNCATE TABLE users CASCADE")
		if err3 != nil {
			log.Fatal("cannot truncate users: ", err3)
		}
	}
	os.Exit(m.Run())
}
