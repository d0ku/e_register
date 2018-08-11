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

func (app *AppContext) mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		userData, err := app.getUserDataFromRequest(w, r)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch userData.UserType {
		case "teacher":
			//display list of schools to choose from.
			app.chooseSchool(w, r)
		case "schoolAdmin":
			//display list of schools to choose from.
			app.chooseSchool(w, r)
		case "student":
			//redirect to /main/student/{student_id}
			//TODO: is this correct http status?
			http.Redirect(w, r, "/main/student/"+userData.UserID, http.StatusSeeOther)
			break
		case "parent":
			//redirect to list of children
			break
		}

		return
	}
}

type chooseSchoolTemplateParse struct {
	Schools  []databasehandling.School
	UserType string
	UserName string
}

func (app *AppContext) chooseSchool(w http.ResponseWriter, r *http.Request) {
	userData, err := app.getUserDataFromRequest(w, r)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	schools, err := app.DbHandler.GetSchoolsDetailsWhereTeacherTeaches(userData.UserID)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//No school to log into.
	if len(schools) == 0 {
		err := app.templates["no_school.gtpl"].Execute(w, nil)

		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else if len(schools) == 1 {
		//Only one school, redirect to it instantly.e
		http.Redirect(w, r, "/main/"+userData.UserType+"/"+strconv.FormatInt(int64(schools[0].Id), 10), http.StatusSeeOther)
	} else {

		//Many schools to choose from.

		templateData := chooseSchoolTemplateParse{schools, userData.UserType, userData.Username}

		err = app.templates["choose_school.gtpl"].Execute(w, templateData)

		if err != nil {
			log.Print(err)
		}
	}
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

//Gets called when client requests /main/teacher/{school_id}
func mainHandleTeacher(w http.ResponseWriter, r *http.Request) {
	fields := strings.Split(r.RequestURI, "/")
	schoolID := fields[len(fields)-1]
	fmt.Println(schoolID)
	w.Write([]byte("lol"))
}
