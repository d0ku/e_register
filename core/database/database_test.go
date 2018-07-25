package database

import (
	"flag"
	"os"
	"testing"
)

var (
	database = flag.Bool("database", false, "run database integration tests")
)

func teardown() {
	DBAddUserTeardown()

}

func TestMain(m *testing.M) {
	flag.Parse()
	if *database {
		temp, err := GetDatabaseHandler("postgres", "test_database")
		if err != nil {
			panic(err)
		}
		dbHandler = temp

		result := m.Run()

		teardown()
		os.Exit(result)
	} else {
		os.Exit(0)
	}
}

func TestDBAddUser(t *testing.T) {
	DBAddUserTeardown() //remove all possible before stuff.
	//TODO: write test.
	DBAddUserTeardown()
}

func DBAddUserTeardown() {
	dbHandler.QueryRow("DELETE FROM Users WHERE username='1234adduser';")
}

func TestDBCheckLoginData(t *testing.T) {
	dbHandler.QueryRow("SELECT add_user('1234logindata','12345678','teacher',1)")
	userData := dbHandler.CheckUserLogin("1234logindata", "12345678", "teacher")

	if userData.exists != true {
		t.Errorf("User exists in database but was not scanned properly.")
	}

	if userData.user_type != "teacher" {
		t.Errorf("User type was not scanned properly")
	}

	if userData.id != 1 {
		t.Errorf("User id was not scanned properly")
	}
}