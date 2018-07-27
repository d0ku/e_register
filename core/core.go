package core

//TODO:	User has to perform action in specified amount of time, add counter on webpage. If he does not, javascript automatically logs him out and his session is deleted from both database and application cache.
import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/d0ku/database_project_go/core/databaseLayer"
	"github.com/d0ku/database_project_go/core/sessions"
)

var dbHandler *databaseLayer.DBHandler

var sessionManager *sessions.SessionManager
var templates map[string]*template.Template

//UserData represents data used to fill in GO HTML templates.
type UserData struct {
	UserName string
	IsLogged bool
}

func getUserDataFromRequest(response http.ResponseWriter, request *http.Request) (*UserData, error) {
	user := &UserData{}
	cookie, err := request.Cookie("sessionID")
	if err != nil {
		return user, errors.New("No cookie set")
	}

	session, err := sessionManager.GetSession(cookie.Value)

	if err != nil {
		//Delete cookie which is notrecognized on server side.
		return user, errors.New("That session does not exist")
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(response, cookie)
	}
	user.UserName = session.Data["username"]
	user.IsLogged = true

	return user, nil
}

func parseAllTemplates(pageFolder string) {
	templates = make(map[string]*template.Template)

	templateFiles, err := ioutil.ReadDir(pageFolder)
	if err != nil {
		panic(err)
	}

	fmt.Println("Parsing templates: ")
	for _, templateFile := range templateFiles {
		if !templateFile.IsDir() {
			name := templateFile.Name()

			//BUG: this regex matches too much at the moment
			regex := regexp.MustCompile("^.*\\.gtpl$")
			//Parse only gtpl files at the moment (html templates).
			if regex.MatchString(name) {
				fmt.Println("---> " + pageFolder + name)
				templates[name] = template.Must(template.ParseFiles(pageFolder + name))
			}
		}
	}
	fmt.Println("Finished parsing.")
}

func mainHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {

		user, err := getUserDataFromRequest(response, request)
		//In theory we don't have to check whether username exists, as parsing template without arguments could just render unlogged site?

		/*
			if err != nil {
				err := templates["index.gtpl"].Execute(response, nil)
				if err != nil {
					log.Fatal(err)
				}

				return
			}
		*/

		err = templates["index.gtpl"].Execute(response, user)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		//POST
		//call some functinos or do something
	}
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("sessionID")
	if err != nil {
		log.Print("Try to log out not yet logged in user.")
		//TODO: display something that user is not even logged in, or just ignore it and redirect to main.
		return
	}

	//Delete cookie and redirect to main.
	cookie.Expires = time.Unix(0, 0)
	http.SetCookie(response, cookie)
	sessionManager.RemoveSession(cookie.Value)

	http.Redirect(response, request, "/main", http.StatusSeeOther)
}

func loginUsers(response http.ResponseWriter, request *http.Request) {
	regex := regexp.MustCompile("/[A-z]*$")
	userType := regex.FindString(request.URL.EscapedPath())[1:]
	//Execute template with correct value to be set as hidden attribute in HTML form.
	err := templates["login_page.gtpl"].Execute(response, userType)
	if err != nil {
		log.Print(err)
	}
}

func loginHandler(response http.ResponseWriter, request *http.Request) {

	fmt.Println(request.Method)
	if request.Method == "GET" {
		//if logged in display personalized site, else display login site

		cookie, err := request.Cookie("sessionID")
		if err != nil {
			//not logged in, cookie does not exist
			err := templates["login.gtpl"].Execute(response, nil)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		session, err := sessionManager.GetSession(cookie.Value)

		if err != nil {
			log.Print("Client side thinks it is logged in, but it is not.")

			//TODO: display something about that he has to relogin (probably redirect should be enough).
			cookie.Expires = time.Unix(0, 0)
			http.SetCookie(response, cookie)
			http.Redirect(response, request, "/main", http.StatusSeeOther)
			//user thinks he is logged in, but it is not true.
			return

		}

		err = templates["login_personal.gtpl"].Execute(response, session.Data["username"])
		if err != nil {
			log.Print(err)
		}
		return

		//logged in

	} else {
		request.ParseForm()
		//Validate data, and check whether it can be used to log into database.
		username := template.HTMLEscapeString(request.Form["username"][0])
		password := template.HTMLEscapeString(request.Form["password"][0])
		userType := template.HTMLEscapeString(request.Form["userType"][0])
		fmt.Println(username)
		fmt.Println(password)

		var checkSchool bool

		if userType == "schoolAdmin" {
			userType = "teacher"
			checkSchool = true
		}
		//TODO: implement js on client-side that checks password length etc.

		user := dbHandler.CheckUserLogin(username, password, userType)

		if !user.Exists {
			fmt.Println("User does not exist!")

		} else {

			if checkSchool {
				schoolID := dbHandler.CheckIfTeacherIsSchoolAdmin(user.Id)
				if schoolID == -1 {
					//ERROR OUT!
					fmt.Println("NO ADMIN!")
					return
				}
				fmt.Println("YUP, ADMIN!")
			}
			fmt.Println("User logged in!")
			//good password and username combination

			//If we have stored sessionID for that user just send it back to him,
			//else create new one.

			sessionID := sessionManager.GetSessionID(username)
			cookie := &http.Cookie{Name: "sessionID", Value: sessionID, MaxAge: 0, Secure: true, HttpOnly: false}
			log.Print("Successfull authentication: " + username)

			/*
				if err != nil {
					//Data was found.
					cookie.Value = sessionID
					log.Print("Successful relogin: " + usernameBase.String)
				} else {
					sessionID := GenerateSessionID(32)
					AddSession(sessionID, username)
					cookie.Value = sessionID
					log.Print("Successful log in: " + usernameBase.String)
				}
			*/

			//Send cookie with sessionID to Client.
			http.SetCookie(response, cookie)

			http.Redirect(response, request, "/main", http.StatusSeeOther)
		}
	}
}

func deleteHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte(request.URL.String()))
}

func registerHandler(response http.ResponseWriter, request *http.Request) {
	user, _ := getUserDataFromRequest(response, request)

	//TODO: write register.gtpl, if user has adequate privileges, (add variable to User Structure) display adding users panel,
	// else display you don't have privileges.
	templates["register.gtpl"].Execute(response, user)
}

//Initialize sets up connection with database, and assigns handlers.
func Initialize(databaseUser string, databaseName string, templatesPath string) {
	//Initialize parsed templates.
	parseAllTemplates(templatesPath)

	//Initialize DB connection.
	//TODO: change user to something more secure (non-root).

	//Could not initialize connection.
	temp, err := databaseLayer.GetDatabaseHandler(databaseUser, databaseName)
	if err != nil {
		panic(err)
	}
	dbHandler = temp

	sessionManager = sessions.GetSessionManager(32, time.Second*60*15)
	//	sessionManager.ReadSessionsFromDatabase()

	//TODO: some kind of login panel, where admin can add new users?
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login/", loginUsers)
	http.Handle("/", http.FileServer(http.Dir("./page/")))
}

//Run starts initialized server on specified port.
func Run(port string) {
	//Certificates are self signed, so they are not worth a penny, if this is supposed to go into production, certificates should be obtained from approppriate organisation.
	fmt.Println()
	fmt.Println("Listen to me at: https://localhost:" + port)

	err := http.ListenAndServeTLS(":"+port, "certs/server.crt", "certs/server.key", nil)

	//Something went wrong with starting HTTPS server.
	if err != nil {
		panic(err)
	}
}
