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

		schools, err := getSchoolsToChoose(session.Data["user_type"], session.Data["id"])

		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if len(schools) == 0 {
			//TODO: write template for this.
			w.Write([]byte("Nie należysz do żadnej szkoły."))
		}

		if len(schools) == 1 {
			http.Redirect(w, r, "/main/"+session.Data["user_type"]+"/"+strconv.FormatInt(int64(schools[0].Id), 10), http.StatusSeeOther)
		}

		//Many schools to choose from.

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

func getSchoolsToChoose(userType string, userID string) ([]databasehandling.School, error) {
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
	Schools  []databasehandling.School
	UserType string
	UserName string
}

/*
func returnSchoolsList(w http.ResponseWriter, r *http.Request) ([]databasehandling.School, error) {
	var err error
	var schools []databasehandling.School

	session, err := getSessionFromRequest(w, r)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	switch session.Data["user_type"] {
	case "teacher":
		schools, err = databasehandling.DbHandler.GetSchoolsDetailsWhereTeacherTeaches(session.Data["id"])
	}

	if err != nil {
		return nil, err
	}

	return schools, nil
}
*/

/*
func handleTeacherChooseSchool(w http.ResponseWriter, r *http.Request) {
	log.Print("TEACHER")
	session, err := getSessionFromRequest(w, r)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Display window to let teacher choose school he wants to see (ha has to teach in them).

	schools, err := databasehandling.DbHandler.GetSchoolsDetailsWhereTeacherTeaches(session.Data["id"])

	if len(schools) == 0 {
		w.Write([]byte("Nie jesteś dodany do żadnej szkoły!"))
		return
	}

	if len(schools) == 1 {
		//TODO: Is this correct status?
		http.Redirect(w, r, "/main/teacher/"+strconv.FormatInt(int64(schools[0].Id), 10), http.StatusSeeOther)
		return
	}

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
*/

func mainHandleTeacher(w http.ResponseWriter, r *http.Request) {
	fields := strings.Split(r.RequestURI, "/")
	schoolID := fields[len(fields)-1]
	fmt.Println(schoolID)
	w.Write([]byte("lol"))
}
