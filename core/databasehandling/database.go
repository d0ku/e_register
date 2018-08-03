package databasehandling

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	//That's the recommended way to do it.
	_ "github.com/lib/pq"
)

var (
	//DbHandler is globally available handler for database, should be only one entry point to DB in whole app.
	DbHandler *DBHandler
)

type DBHandler struct {
	*sql.DB
}

func GetDatabaseHandler(username string, dbName string, dbPassword string, sslmode string) (*DBHandler, error) {
	var err error
	//TODO: connect to postgres by SSL (sslmode=verify-full)

	connStr := "user=" + username + " dbname=" + dbName + " sslmode=" + sslmode + " password=" + dbPassword
	dbConnection, err := sql.Open("postgres", connStr)

	if err == nil {
		//Check if connection can be established (password matches etc.)
		err = dbConnection.Ping()
	}

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

func (handler *DBHandler) CheckIfTeacherIsSchoolAdmin(teacherID int) int {
	fmt.Println(teacherID)
	query := "SELECT * FROM check_if_teacher_is_school_admin(" + strconv.Itoa(teacherID) + ");"
	fmt.Println(query)
	var output = -1
	err := handler.QueryRow(query).Scan(&output)

	if err != nil {
		log.Print(err)
	}

	return output
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
