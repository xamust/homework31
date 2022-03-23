package store_test

import (
	"os"
	"testing"
)

var dataBaseURL string

func TestMain(m *testing.M) {
	dataBaseURL = os.Getenv("DATABASE_URL")
	if dataBaseURL == "" {
		dataBaseURL = "dbname=testDB user=postgres password=example host=localhost port=5433 sslmode=disable"

	}
	os.Exit(m.Run())
}
