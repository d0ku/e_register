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

	"github.com/d0ku/e_register/core/logging"
	"github.com/d0ku/e_register/core/sessions"
)

var (
	sessionManager *sessions.SessionManager
	//templates contains parsed HTML templates.
	templates map[string]*template.Template

	//ErrNoSuchSession is returned when request contains sessionID cookie, but sessionManager can't find it.
	ErrNoSuchSession = errors.New("Session with such sessionID does not exist")
	//ErrNoSuchSessionCookie is returned when request does not contain sessionID cookie.
	ErrNoSuchSessionCookie = errors.New("No sessionID cookie set")
)

//Gets session data from request and automatically handles:
//- no cookie at all
//- incorrect cookie
//and deletes cookie from client side in that case.
func getSessionFromRequest(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return nil, ErrNoSuchSessionCookie
	}

	session, err := sessionManager.GetSession(cookie.Value)

	if err != nil {
		//Delete cookie which is not recognized on server side.
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
		return nil, ErrNoSuchSession
	}

	return session, nil
}

func parseAllTemplates(pageFolder string) {
	templates = make(map[string]*template.Template)

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
}

func redirectToLoginPageIfUserNotLogged(h http.Handler, messagePorts ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := getSessionFromRequest(w, r)

		if err != nil {
			log.Print("LOGIN|User from: " + r.RemoteAddr + " tried to access app page without privileges.")
			templates["not_logged.gtpl"].Execute(w, nil)
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

//Initialize assigns handlers to provided mux and sets up sessionManager with provided session life time.
//It also parses all HTML templates located under templatesPath.
func Initialize(templatesPath string, cookieLifeTime time.Duration, mux *logging.MuxController) {
	//Parse all HTML templates from provided directory.
	parseAllTemplates(templatesPath)

	//Initialize session manager.
	sessionManager = sessions.GetSessionManager(32, cookieLifeTime)

	//Initialize timeouts after too many login tries module.
	loginController := sessions.GetLoginTriesController()

	//Create fileServer to deliver static content.
	fileServer := http.StripPrefix("/page/", http.FileServer(http.Dir("./page/server_root/")))

	mux.Handle("/login", loginHandlerDecorator(cookieLifeTime, loginController))
	mux.Handle("/logout", redirectToLoginPageIfUserNotLogged(http.HandlerFunc(logoutHandler)))
	mux.Handle("/main", redirectToLoginPageIfUserNotLogged(http.HandlerFunc(mainHandler)))
	mux.Handle("/login/", http.HandlerFunc(loginUsers))
	mux.Handle("/", http.HandlerFunc(redirectToLogin))
	mux.Handle("/page/", fileServer)
}
