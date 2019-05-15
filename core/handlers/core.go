package handlers

//TODO:	User has to perform action in specified amount of time, add counter of cookie lifetime on webpage.

//User session life period is stored in cookie and removed after time ends.

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/d0ku/e_register/core/databasehandling"
	"github.com/d0ku/e_register/core/logging"
	"github.com/d0ku/e_register/core/sessions"
)

//AppContext defines data that all handlers have access to.
type AppContext struct {
	sessionManager       sessions.SessionManager
	templates            map[string]*template.Template
	DbHandler            databasehandling.DBHandler
	cookieLifeTime       time.Duration
	loginTriesController *sessions.LoginTriesController
}

var (
	//ErrNoSuchSession is returned when request contains sessionID cookie, but sessionManager can't find it.
	ErrNoSuchSession = errors.New("Session with such sessionID does not exist")
	//ErrNoSuchSessionCookie is returned when request does not contain sessionID cookie.
	ErrNoSuchSessionCookie = errors.New("No sessionID cookie set")
)

//Gets session data from request and automatically handles:
//- no cookie at all
//- incorrect cookie
//and deletes cookie from client side in that case.
func (app *AppContext) getSessionFromRequest(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return nil, ErrNoSuchSessionCookie
	}

	session, err := app.sessionManager.GetSession(cookie.Value)

	if err != nil {
		//Delete cookie which is not recognized on server side.
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
		return nil, ErrNoSuchSession
	}

	return session, nil
}

//Gets user data from request and automatically handles:
//- no cookie at all
//- incorrect cookie
//and deletes cookie from client side in that case.
func (app *AppContext) getUserDataFromRequest(w http.ResponseWriter, r *http.Request) (*sessions.UserData, error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return nil, ErrNoSuchSessionCookie
	}

	data, err := app.sessionManager.GetUserData(cookie.Value)

	if err != nil {
		//Delete cookie which is not recognized on server side.
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
		return nil, ErrNoSuchSession
	}

	return data, nil
}

func parseAllTemplates(pageFolder string) map[string]*template.Template {
	templates := make(map[string]*template.Template)

	templateFiles, err := ioutil.ReadDir(pageFolder)
	if err != nil {
		panic(err)
	}

	log.Print("TEMPLATES|Parsing started...")
	for _, templateFile := range templateFiles {
		if !templateFile.IsDir() {
			name := templateFile.Name()

			regex := regexp.MustCompile("^.*\\.gtpl$")
			//Parse only gtpl files at the moment (html templates).

			if regex.MatchString(name) {
				log.Print("TEMPLATES|Parsing ---> " + pageFolder + name)
				templates[name] = template.Must(template.ParseFiles(pageFolder + name))
			}
		}
	}
	log.Print("TEMPLATES|Parsing finished...")

	return templates
}

func (app *AppContext) checkUserLogon(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.getSessionFromRequest(w, r)

		if err != nil {
			log.Print("LOGIN|User from: " + r.RemoteAddr + " tried to access app page without privileges.")
			err := app.templates["not_logged.gtpl"].Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
		h.ServeHTTP(w, r)
	})
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (app *AppContext) checkPermission(h http.Handler, userTypeIn string) http.Handler {
	//TODO: that checker is really simple, add more stuff.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, err := app.getUserDataFromRequest(w, r)
		if err != nil {
			log.Print(err)
			err := app.templates["not_logged.gtpl"].Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		if userData.UserType != userTypeIn {
			err := app.templates["no_permission.gtpl"].Execute(w, userTypeIn)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		h.ServeHTTP(w, r)
	})
}

//Initialize assigns handlers to provided mux and sets up sessionManager with provided session life time.
//It also parses all HTML templates located under templatesPath.
func Initialize(templatesPath string, cookieLifeTime time.Duration, mux *logging.MuxController, db databasehandling.DBHandler) {
	//Parse all HTML templates from provided directory.
	templates := parseAllTemplates(templatesPath)

	//Initialize session manager.
	sessionManager := sessions.GetSessionManager(32, cookieLifeTime)

	//Initialize timeouts after too many login tries module.
	loginController := sessions.GetLoginTriesController()

	//Create fileServer to deliver static content.
	fileServer := http.StripPrefix("/page/", http.FileServer(http.Dir("./page/server_root/")))

	//TODO: pass arguments in different way
	appContext := &AppContext{sessionManager, templates, db, cookieLifeTime, loginController}

	///main/{user_type}/{school_id}
	mux.Handle("/login", http.HandlerFunc(appContext.loginHandler))
	mux.Handle("/logout", appContext.checkUserLogon(http.HandlerFunc(appContext.logoutHandler)))
	mux.Handle("/main/", appContext.checkUserLogon(http.HandlerFunc(appContext.mainHandler)))
	mux.Handle("/main/teacher/", appContext.checkPermission(appContext.checkUserLogon(http.HandlerFunc(appContext.mainHandleTeacher)), "teacher"))
	mux.Handle("/main/student/", appContext.checkPermission(appContext.checkUserLogon(http.HandlerFunc(mainHandleStudent)), "student"))
	mux.Handle("/main/schoolAdmin/", appContext.checkPermission(appContext.checkUserLogon(http.HandlerFunc(mainHandleSchoolAdmin)), "schoolAdmin"))
	mux.Handle("/main/parent/", appContext.checkPermission(appContext.checkUserLogon(http.HandlerFunc(mainHandleParent)), "parent"))
	mux.Handle("/login/", http.HandlerFunc(appContext.loginUsers))
	mux.Handle("/", http.HandlerFunc(redirectToLogin))
	mux.Handle("/page/", fileServer)
}
