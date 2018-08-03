package handlers

//All requests that are called in this file should already pass through redirectToLogin handler decorator, so we can be sure that they have cookie with correct value set up.

import (
	"log"
	"net/http"
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

func mainHandleTeacher(w http.ResponseWriter, r *http.Request) {
	log.Print("TEACHER")
	session, err := getSessionFromRequest(w, r)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates["index.gtpl"].Execute(w, session.Data["username"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
	}
}
