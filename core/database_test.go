package core

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	temp, err := GetDatabaseHandler("postgres", "test_database")
	if err != nil {
		panic(err)
	}
	dbHandler = temp

	os.Exit(m.Run())
}

func TestDBAddUser(t *testing.T) {
	dbHandler.QueryRow("DELETE FROM Users WHERE username='1234adduser'")
	//TODO: write test.
}

func TestDBCheckLoginData(t *testing.T) {
	dbHandler.QueryRow("SELECT add_user('1234logindata','12345678','teacher',1)")
	userData := dbHandler.CheckUserLogin("1234logindata", "12345678")

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
