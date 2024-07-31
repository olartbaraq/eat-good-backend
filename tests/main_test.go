package all_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/GoogleCloudPlatform/golang-samples/run/helloworld/db/sqlc"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/utils"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {
	// This function is to perform the main test.

	config, err := utils.LoadDBConfig("..")
	if err != nil {
		log.Fatal("Could not load env config", err)
	}

	conn, err := sql.Open(config.DBdriver, config.DBsource)
	if err != nil {
		log.Fatalf("There was an error connecting to database: %v", err)
	}
	testQueries = db.New(conn)
	os.Exit(m.Run())
}
