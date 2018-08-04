package handlers

//All requests that are called in this file should already pass through redirectToLogin handler decorator, so we can be sure that they have cookie with correct value set up.

import (
	"log"
	"net/http"

	"github.com/d0ku/e_register/core/databasehandling"
)

/*
func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, err := getSessionFromRequest(w, r)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

	}
}
*/

func mainHandleSchoolAdmin(w http.ResponseWriter, r *http.Request) {
	log.Print("ADMIN")

}

func mainHandleParent(w http.ResponseWriter, r *http.Request) {
	log.Print("PARENT")

}

func mainHandleStudent(w http.ResponseWriter, r *http.Request) {
	log.Print("STUDENT")

}

type chooseSchoolTemplateParse struct {
	Schools  []*databasehandling.School
	UserType string
	UserName string
}

func mainHandleTeacher(w http.ResponseWriter, r *http.Request) {
	log.Print("TEACHER")
	session, err := getSessionFromRequest(w, r)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Display window to let teacher choose school he wants to see (ha has to teach in them).

	schools, err := databasehandling.DbHandler.GetSchoolsDetailsWhereTeacherTeaches(session.Data["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := chooseSchoolTemplateParse{schools, session.Data["user_type"], session.Data["username"]}
	err = templates["choose_school.gtpl"].Execute(w, templateData)

	if err != nil {
		log.Print(err)
	}
}
