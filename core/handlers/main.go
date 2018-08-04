package handlers

//All requests that are called in this file should already pass through redirectToLogin handler decorator, so we can be sure that they have cookie with correct value set up.

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/d0ku/e_register/core/databasehandling"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, err := getSessionFromRequest(w, r)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		schools, err := getDataToChoose(session.Data["user_type"], session.Data["id"])

		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//No school/children to log into.
		if len(schools) == 0 {
			err := templates["no_school.gtpl"].Execute(w, nil)

			if err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if len(schools) == 1 {
			//Only one school or children, redirect to it instantly.e
			http.Redirect(w, r, "/main/"+session.Data["user_type"]+"/"+strconv.FormatInt(int64(schools[0].Id), 10), http.StatusSeeOther)
		} else {

			//Many schools or children to choose from.

			templateData := chooseSchoolTemplateParse{schools, session.Data["user_type"], session.Data["username"]}

			err = templates["choose_school.gtpl"].Execute(w, templateData)

			if err != nil {
				log.Print(err)
			}

			/*
				switch session.Data["user_type"] {

				case "teacher":
					mainHandleTeacher(w, r)
				case "schoolAdmin":
					mainHandleSchoolAdmin(w, r)
				case "student":
					mainHandleStudent(w, r)
				case "parent":
					mainHandleParent(w, r)

				default:
					//TODO: should such unknown user be automatically logged out and his cookie deleted?
					log.Print("BUG|There was a try to log in as user of not specified type, this should not happen.")
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			*/

		}
	}
}

//getDataToChoose returns appropriate data basing on provided userType and userID.
//It errors out only when there was some error in databasehandling.
//It returns empty array if no data is associated with user.
func getDataToChoose(userType string, userID string) ([]databasehandling.School, error) {
	//Return schools in case of teachers and admins, children in case of parents and class in case of student.
	var schools []databasehandling.School
	var err error

	switch userType {
	case "teacher":
		schools, err = databasehandling.DbHandler.GetSchoolsDetailsWhereTeacherTeaches(userID)
	case "schoolAdmin":

	case "student":

	case "parent":

	}

	return schools, err
}

//Gets called when client requests /main/schoolAdmin/{school_id}
func mainHandleSchoolAdmin(w http.ResponseWriter, r *http.Request) {
	log.Print("ADMIN")

}

//Gets called when client requests /main/parent/{student_id} where {student_id} is id of one of its child.
func mainHandleParent(w http.ResponseWriter, r *http.Request) {
	log.Print("PARENT")

}

//Gets called when client requests /main/parent/{student_id}
func mainHandleStudent(w http.ResponseWriter, r *http.Request) {
	log.Print("STUDENT")

}

type chooseSchoolTemplateParse struct {
	Schools  []databasehandling.School
	UserType string
	UserName string
}

//Gets called when client requests /main/teacher/{school_id}
func mainHandleTeacher(w http.ResponseWriter, r *http.Request) {
	fields := strings.Split(r.RequestURI, "/")
	schoolID := fields[len(fields)-1]
	fmt.Println(schoolID)
	w.Write([]byte("lol"))
}
