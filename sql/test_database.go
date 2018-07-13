package sql_test

import (
	"database/sql"
	"os"
	"testing"
)

var dbConnection *sql.DB

func TestUserAdding(t *testing.T) {
	var returnValue bool
	commands := [...]string{
		"SELECT add_user('testUser','testPassword','teacher',1);",
		"SELECT add_user('testUser','testPassword','teacher',1); --should return false (duplicate)",
	}

	returnValues := [...]bool{
		true,
		false,
	}
	for index, command := range commands {

		dbConnection.QueryRow(command).Scan(&returnValue)
		if returnValue != returnValues[index] {
			t.Errorf(command)
		}
	}
}

//TestMain sets up databaseConnection and runs tests.
func TestMain(m *testing.M) {
	//TODO: run with SSL.
	var err error
	connStr := "user=postgres dbname=test_database sslmode=disable"
	dbConnection, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
