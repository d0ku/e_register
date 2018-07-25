package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBHandler struct {
	*sql.DB
}

func GetDatabaseHandler(username string, dbName string) (*DBHandler, error) {
	var err error
	//TODO: connect to postgres by SSL (sslmode=verify-full)

	connStr := "user=" + username + " dbname=" + dbName + " sslmode=disable"
	dbConnection, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	return &DBHandler{dbConnection}, nil
}

type UserLoginData struct {
	Exists    bool
	User_type string
	Id        int
}

func (handler *DBHandler) CheckUserLogin(username string, password string, userType string) *UserLoginData {
	query := "SELECT * FROM check_login_data('" + username + "','" + password + "','" + userType + "');"
	fmt.Println(query)
	var exists bool
	var user_type string
	var id int
	err := handler.QueryRow(query).Scan(&exists, &user_type, &id)
	if err != nil {
		log.Print(err)
	}
	return &UserLoginData{exists, user_type, id}
}

func (handler *DBHandler) AddSession(session_id string, username string) bool {
	var result bool
	query := "SELECT add_session('" + session_id + "','" + username + "');"
	fmt.Println(query)
	err := handler.QueryRow(query).Scan(&result)
	if err != nil {
		log.Print(err)
	}
	return result
}

func (handler *DBHandler) DeleteSession(session_id string) {
	query := "SELECT delete_session('" + session_id + "');"
	fmt.Println(query)
	err := handler.QueryRow(query)
	if err != nil {
		log.Print(err)
	}
}
