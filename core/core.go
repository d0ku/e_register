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
	"math"
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

//Gets session data from request and automatically handles:
//- no cookie at all
//- incorrect cookie
//and deletes cookie from client side in that case.
func getSessionFromRequest(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return nil, errors.New("No cookie set")
	}

	session, err := sessionManager.GetSession(cookie.Value)

	if err != nil {
		//Delete cookie which is not recognized on server side.
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
		return nil, errors.New("That session does not exist")
	}

	return session, nil
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
	log.Print(r.Method)
	if r.Method == "GET" {
		var err error
		//We know that request has been checked previously so there is no need to check for error.
		session, _ := getSessionFromRequest(w, r)

		err = templates["index.gtpl"].Execute(w, session.Data["username"])
		if err != nil {
			log.Print(err)
		}

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

	//TODO: display some info about successful log out?
	http.Redirect(w, r, "/login", http.StatusFound)
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

func loginHandlerGET(w http.ResponseWriter, r *http.Request) {
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

	_, err = sessionManager.GetSession(cookie.Value)

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

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

type UserLoginTry struct {
	UserType   string
	UserName   string
	Timeout    int
	HasTimeout bool
}

func loginHandlerDecorator(cookieLifeTime time.Duration, loginTriesController *LoginTriesController) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			loginHandlerGET(w, r)
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

			userLoginTry := &UserLoginTry{UserType: userType, UserName: username, Timeout: 0, HasTimeout: false}
			fmt.Println(userLoginTry)
			if userType == "schoolAdmin" {
				userType = "teacher"
				checkSchool = true
			}

			user := dbHandler.CheckUserLogin(username, password, userType)

			if !user.Exists {
				loginTriesController.AddTry(r.RemoteAddr)
				timeoutSecs := loginTriesController.GetTimeoutLeft(r.RemoteAddr)
				if timeoutSecs != 0 {
					userLoginTry.Timeout = timeoutSecs
					userLoginTry.HasTimeout = true
				}

				//TODO: add javascript count down in template file, so user could see when his timeout is done.
				log.Print("Unsuccessful try to log in from:" + r.RemoteAddr)
				err := templates["login_error.gtpl"].Execute(w, userLoginTry)
				if err != nil {
					log.Print(err)
				}
			} else {
				if checkSchool {
					schoolID := dbHandler.CheckIfTeacherIsSchoolAdmin(user.Id)
					if schoolID == -1 {
						//There is no schoolAdmin with such id.

						//Timeouts are also issued when someone tries to log in as admin, when user is only a teacher.
						loginTriesController.AddTry(r.RemoteAddr)
						timeoutSecs := loginTriesController.GetTimeoutLeft(r.RemoteAddr)
						if timeoutSecs != 0 {
							userLoginTry.Timeout = timeoutSecs
							userLoginTry.HasTimeout = true
						}

						//TODO: add javascript count down in template file, so user could see when his timeout is done.
						//When teacher tries to login as admin, he gets same error message as if his username and password didn't match. That's true after all.
						err := templates["login_error.gtpl"].Execute(w, userLoginTry)
						if err != nil {
							log.Print(err)
						}

						log.Print("Unsuccessful try to log in as admin (no permissions) from:" + username)
						return
					}
					log.Print("Successful admin logon from:" + r.RemoteAddr)
				}
				loginTriesController.ResetTries(r.RemoteAddr)
				//TODO: Is that kind of logging neccessary? GDPR compliance and so on?
				log.Print("Logon as: " + username + " from:" + r.RemoteAddr)

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

//TODO: consider moving below code to another source file.
//Is this memory safe to never purge addresses which didn't login successfully?

//LoginTriesController enables application to easily control how many times user tried to log in and give him specific timeouts.
type LoginTriesController struct {
	tries      map[string]int
	timeoutEnd map[string]time.Time
}

func (controller *LoginTriesController) setTimeout(origin string) {
	var addTime time.Duration
	tries := controller.tries[origin]

	//TODO: define better policy for timeouts.
	if tries >= 100 {
		addTime = time.Second * 30
	} else if tries >= 50 {
		addTime = time.Second * 20
	} else if tries >= 10 {
		addTime = time.Second * 15
	} else if tries >= 5 {
		addTime = time.Second * 10
	}

	controller.timeoutEnd[origin] = time.Now().Add(addTime)
}

//AddTry add one try to user, if tries count exceeds specified value, it calls setTimeout function which set timeouts according to specified policy.
func (controller *LoginTriesController) AddTry(origin string) {
	_, ok := controller.timeoutEnd[origin]
	if ok {
		//Don't do anything if there is a timeout already.
		return
	}
	_, ok = controller.tries[origin]

	if !ok {
		controller.tries[origin] = 0
	}

	controller.tries[origin]++

	if controller.tries[origin] >= 5 {
		controller.setTimeout(origin)
	}
}

//ResetTries resets try counter.
func (controller *LoginTriesController) ResetTries(origin string) {
	delete(controller.tries, origin)
}

//GetTimeoutLeft returns 0 if there is no timeout left, or int (seconds) representing how long user has to wait.
//If timeout left is equal to 0, user can try to log in.
//Else function should return timeout left (in seconds).
func (controller *LoginTriesController) GetTimeoutLeft(origin string) int {
	timeout, ok := controller.timeoutEnd[origin]

	if !ok {
		//There is no timeout set.
		return 0
	}

	if time.Now().After(timeout) {
		//Timeout already passed, can be deleted.
		delete(controller.timeoutEnd, origin)
		return 0
	}

	//Return time left.
	timeLeft := int(math.Round(timeout.Sub(time.Now()).Seconds()))
	if timeLeft == 0 {
		delete(controller.timeoutEnd, origin)
		return 0
	}

	return timeLeft
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.String()))
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

func redirectWithErrorToLogin(h http.Handler, messagePorts ...string) func(http.ResponseWriter, *http.Request) {
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
			return
		}

		h.ServeHTTP(w, r)
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

	loginController := &LoginTriesController{make(map[string]int), make(map[string]time.Time)}

	http.HandleFunc("/login", loginHandlerDecorator(cookieLifeTime, loginController))
	http.HandleFunc("/logout", redirectWithErrorToLogin(http.HandlerFunc(logoutHandler)))
	//http.HandleFunc("/delete", redirectWithErrorToLogin(deleteHandler))
	http.HandleFunc("/main", redirectWithErrorToLogin(http.HandlerFunc(mainHandler)))
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
