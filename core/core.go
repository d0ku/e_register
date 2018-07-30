package core

//TODO:	User has to perform action in specified amount of time, add counter of cookie lifetime on webpage.

//User session life period is stored in cookie and removed after time ends.

//All rs are checked for sessionID cookie when we get them, so there is no need to check for errors in getting that cookie in later rs.
import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
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

func getUserDataFromRequest(w http.ResponseWriter, r *http.Request) (*UserData, error) {
	user := &UserData{}
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return user, errors.New("No cookie set")
	}

	session, err := sessionManager.GetSession(cookie.Value)

	if err != nil {
		//Delete cookie which is notrecognized on server side.
		return user, errors.New("That session does not exist")
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
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

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		user, err := getUserDataFromRequest(w, r)
		//In theory we don't have to check whether username exists, as parsing template without arguments could just render unlogged site?

		/*
			if err != nil {
				err := templates["index.gtpl"].Execute(w, nil)
				if err != nil {
					log.Fatal(err)
				}

				return
			}
		*/

		err = templates["index.gtpl"].Execute(w, user)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		//POST
		//call some functinos or do something
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		log.Print("Try to log out not yet logged in user.")
		//TODO: display something that user is not even logged in, or just ignore it and redirect to main.
		return
	}

	//Delete cookie and redirect to main.
	cookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, cookie)
	sessionManager.RemoveSession(cookie.Value)

	http.Redirect(w, r, "/main", http.StatusFound)
}

func loginUsers(w http.ResponseWriter, r *http.Request) {
	fields := strings.Split(r.RequestURI, "/")
	userType := fields[len(fields)-1]

	//Execute template with correct value to be set as hidden attribute in HTML form.
	err := templates["login_page.gtpl"].Execute(w, userType)
	if err != nil {
		log.Print(err)
	}
}

func loginHandlerDecorator(cookieLifeTime time.Duration) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			cookie, err := r.Cookie("sessionID")
			if err != nil {
				log.Print("Normal try to log in from: " + r.RemoteAddr)
				//User is not logged in, cookie does not exist, normal use-case.

				err := templates["login.gtpl"].Execute(w, nil)
				if err != nil {
					log.Print(err)
				}
				return
			}

			session, err := sessionManager.GetSession(cookie.Value)

			if err != nil {
				log.Print("Incorrect cookie on user side.")
				//User tried to log in with expired cookie or he is trying to do something malicious.

				//Remove expired cookie from client side.
				cookie.Expires = time.Unix(0, 0)
				http.SetCookie(w, cookie)

				//Show user info that he is not logged in.
				templates["not_logged.gtpl"].Execute(w, nil)
				return
			}

			//If logged user tries to access /login page, we redirect him to /main.
			//BUG: [possible] Is this correct http status for such case?

			http.Redirect(w, r, "/main", http.StatusFound)
			err = templates["login_personal.gtpl"].Execute(w, session.Data["username"])
			if err != nil {
				log.Print(err)
			}

		} else { //POST request.
			r.ParseForm()
			//Validate data, and check whether it can be used to log into database.
			username := template.HTMLEscapeString(r.Form["username"][0])
			password := template.HTMLEscapeString(r.Form["password"][0])
			userType := template.HTMLEscapeString(r.Form["userType"][0])

			//DEBUG
			fmt.Println(username)
			fmt.Println(password)
			//END OF DEBUG

			var checkSchool bool

			if userType == "schoolAdmin" {
				userType = "teacher"
				checkSchool = true
			}

			user := dbHandler.CheckUserLogin(username, password, userType)

			if !user.Exists {
				//TODO: implement timeouts dependind on number of tries from address.
				log.Print("Unsuccessful try to log in from:" + r.RemoteAddr)
			} else {
				if checkSchool {
					schoolID := dbHandler.CheckIfTeacherIsSchoolAdmin(user.Id)
					if schoolID == -1 {
						//There is no schoolAdmin with such id.
						log.Print("Try to log in as admin (no permissions): " + username)
					}
					log.Print("Successful admin logon from:" + r.RemoteAddr)
				}
				//TODO: Is that kind of logging neccessary?
				log.Print("Logon as: " + username + " from " + r.RemoteAddr)

				//We always create new session for users who don't have valid cookies.
				sessionID := sessionManager.GetSessionID(username)

				//Send cookie with defined expiration time and sessionID value to user.
				cookie := &http.Cookie{Name: "sessionID", Value: sessionID, Expires: time.Now().Add(cookieLifeTime * time.Second), Secure: true, HttpOnly: false}

				//Send cookie with sessionID to Client.
				http.SetCookie(w, cookie)

				//Redirect user to main.
				//TODO: display some info about successfull login?
				//TODO: is that correct http status?

				http.Redirect(w, r, "/main", http.StatusSeeOther)
			}
		}
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.String()))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := getUserDataFromRequest(w, r)

	//TODO: write register.gtpl, if user has adequate privileges, (add variable to User Structure) display adding users panel,
	// else display you don't have privileges.
	templates["register.gtpl"].Execute(w, user)
}

func redirectToHTTPS(h http.Handler, ports ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		host := r.Host
		if i := strings.Index(host, ":"); i != -1 {
			host = host[:i]
		}
		redirectAddress := "https://" + host + ":" + ports[0] + r.RequestURI

		log.Print(redirectAddress)
		http.Redirect(w, r, redirectAddress, http.StatusMovedPermanently)
	})
}

func redirectWithErrorToLogin(h func(http.ResponseWriter, *http.Request), messagePorts ...string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionID")

		if err != nil {
			log.Print(err)
			templates["not_logged.gtpl"].Execute(w, nil)
			return
		}

		_, ok := sessionManager.GetSession(cookie.Value)

		if ok != nil {
			log.Print("User from: " + r.RemoteAddr + " tried to log in with incorrect cookie.")
			templates["not_logged.gtpl"].Execute(w, nil)
		}
	}
}

func placeHolderHandler(w http.ResponseWriter, r *http.Request) {
}

//Initialize sets up connection with database, and assigns handlers.
func Initialize(databaseUser string, databaseName string, templatesPath string, cookieLifeTime time.Duration) {
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

	http.HandleFunc("/login", loginHandlerDecorator(cookieLifeTime))
	http.HandleFunc("/logout", redirectWithErrorToLogin(logoutHandler))
	//http.HandleFunc("/delete", redirectWithErrorToLogin(deleteHandler))
	http.HandleFunc("/main", redirectWithErrorToLogin(mainHandler))
	//	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login/", loginUsers)
	http.Handle("/", http.FileServer(http.Dir("./page/")))
}

//RunTLS starts initialized server on specified port with TLS.
func RunTLS(HTTPSport string, HTTPort string, redirectHTTPtoHTTPS bool, hostname string, serverCert string, serverKey string) {
	//Certificates are self signed, so they are not worth a penny, if this is supposed to go into production, certificates should be obtained from approppriate organisation.

	if redirectHTTPtoHTTPS {
		go func() {
			err := http.ListenAndServe(":"+HTTPort, redirectToHTTPS(http.HandlerFunc(placeHolderHandler), HTTPSport))
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	fmt.Println()
	fmt.Println("Listen to me at: https://" + hostname + ":" + HTTPSport)

	err := http.ListenAndServeTLS(":"+HTTPSport, serverCert, serverKey, nil)

	//Something went wrong with starting HTTPS server.
	if err != nil {
		panic(err)
	}
}
