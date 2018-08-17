package databasehandling

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	//That's the recommended way to do it.
	_ "github.com/lib/pq"
)

var (
	//ErrCouldNotGetRows is returned when there is a problem with querying
	ErrCouldNotGetRows = errors.New("Could not get rows")
)

//DBHandler defines default interactions with db.
type DBHandler interface {
	CheckUserLogin(string, string, string) *UserLoginData
	CheckIfTeacherIsSchoolAdmin(int) int
	GetSchoolsDetailsWhereTeacherTeaches(string) ([]School, error)
}

type dbHandlerStruct struct {
	*sql.DB
	statements map[string]*sql.Stmt
}

//GetDatabaseHandler returns handler compliant with specified options.
//TODO: implement specyfing remote server address for SQL connection.
func GetDatabaseHandler(username string, dbName string, dbPassword string, sslmode string) (DBHandler, error) {
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

	statements := getStatements(dbConnection)
	handler := &dbHandlerStruct{dbConnection, statements}

	return DBHandler(handler), nil
}

//TODO: Is this really worth it?
func getStatements(dbConn *sql.DB) map[string]*sql.Stmt {
	statements := make(map[string]*sql.Stmt)

	temp, err := dbConn.Prepare("SELECT * FROM check_login_data($1,$2,$3);")
	if err != nil {
		panic(err)
	}
	statements["login"] = temp

	temp, err = dbConn.Prepare("SELECT * FROM check_if_teacher_is_school_admin($1);")
	if err != nil {
		panic(err)
	}
	statements["check_school_admin"] = temp

	temp, err = dbConn.Prepare("SELECT * FROM get_schools_details_where_teacher_teaches($1);")
	if err != nil {
		panic(err)
	}
	statements["check_schools_details"] = temp

	return statements
}

func (handler *dbHandlerStruct) CheckUserLogin(username string, password string, userType string) *UserLoginData {
	var exists bool
	var userTypeOut string
	var id int
	var change_password bool
	//err = handler.QueryRow(query).Scan(&exists, &userTypeOut, &id, &change_password)
	err := handler.statements["login"].QueryRow(username, password, userType).Scan(&exists, &userTypeOut, &id, &change_password)
	if err != nil {
		log.Print(err)
	}
	return &UserLoginData{exists, userTypeOut, id, change_password}
}

func (handler *dbHandlerStruct) CheckIfTeacherIsSchoolAdmin(teacherID int) int {
	var output = -1
	err := handler.statements["check_school_admin"].QueryRow(teacherID).Scan(&output)

	if err != nil {
		log.Print(err)
	}

	return output
}

func (handler *dbHandlerStruct) GetSchoolsDetailsWhereTeacherTeaches(teacherID string) ([]School, error) {
	//We are not returning pointers, because Go templates could not handle them at the moment.

	rows, err := handler.statements["check_schools_details"].Query(teacherID)
	if err != nil {
		log.Print(err)
		return nil, ErrCouldNotGetRows
	}

	schools := make([]School, 0)

	for rows.Next() {
		var id int
		var fullName string
		var city string
		var street string
		var schoolType string

		err := rows.Scan(&id, &fullName, &city, &street, &schoolType)
		if err != nil {
			log.Print(err)
		}
		school := School{Id: id, FullName: fullName, City: city, Street: street, SchoolType: schoolType}
		//DEBUG
		fmt.Println(school)
		schools = append(schools, school)
	}

	return schools, nil
}
