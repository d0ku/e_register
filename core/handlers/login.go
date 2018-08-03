package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/d0ku/e_register/core/databasehandling"
	"github.com/d0ku/e_register/core/sessions"
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")

	if err != nil {
		//No cookie found, nothing to do except redirect.
	} else {

		session, err := sessionManager.GetSession(cookie.Value)
		if err != nil {
			log.Print(err)
		} else {
			log.Print("LOGIN|Successfully logged out: " + session.Data["username"] + " from:" + r.RemoteAddr)
		}

		//Delete cookie and redirect to main.
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
		sessionManager.RemoveSession(cookie.Value)
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func loginUsers(w http.ResponseWriter, r *http.Request) {
	fields := strings.Split(r.RequestURI, "/")
	userType := fields[len(fields)-1]

	//Execute template with correct value to be set as hidden attribute in HTML form.
	err := templates["login_form.gtpl"].Execute(w, userType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
	}
}

func loginHandlerGET(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		//User is not logged in, cookie does not exist, normal use-case.

		err := templates["login.gtpl"].Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err)
		}
		return
	}

	_, err = sessionManager.GetSession(cookie.Value)

	if err != nil {
		//User tried to log in with expired cookie or he is trying to do something malicious.
		log.Print("LOGIN|Incorrect cookie on user side.")

		//Remove expired cookie from client side.
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)

		//Show user information page saying that he is not logged in.
		templates["not_logged.gtpl"].Execute(w, nil)
		return
	}

	//If logged user tries to access /login page, we redirect him to /main.

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

func loginHandlerDecorator(cookieLifeTime time.Duration, loginTriesController *sessions.LoginTriesController) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			loginHandlerGET(w, r)
		} else { //POST request.
			r.ParseForm()
			//Validate data, and check whether it can be used to log into database.
			username := template.HTMLEscapeString(r.Form["username"][0])
			password := template.HTMLEscapeString(r.Form["password"][0])
			userType := template.HTMLEscapeString(r.Form["userType"][0])

			var checkSchool bool

			userLoginTry := &sessions.UserLoginTry{UserType: userType, UserName: username, Timeout: 0, HasTimeout: false}
			//log.Println(userLoginTry)

			if userType == "schoolAdmin" {
				userType = "teacher"
				checkSchool = true
			}

			user := databasehandling.DbHandler.CheckUserLogin(username, password, userType)

			if !user.Exists {
				loginTriesController.AddTry(r.RemoteAddr)
				timeoutSecs := loginTriesController.GetTimeoutLeft(r.RemoteAddr)
				if timeoutSecs != 0 {
					userLoginTry.Timeout = timeoutSecs
					userLoginTry.HasTimeout = true
				}

				log.Print("LOGIN|Unsuccessful try to log in from:" + r.RemoteAddr)
				err := templates["login_error.gtpl"].Execute(w, userLoginTry)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					log.Print(err)
				}

			} else {
				if checkSchool {
					schoolID := databasehandling.DbHandler.CheckIfTeacherIsSchoolAdmin(user.Id)
					if schoolID == -1 {
						//There is no schoolAdmin with such id.

						//Timeouts are also issued when someone tries to log in as admin, when user is only a teacher.
						loginTriesController.AddTry(r.RemoteAddr)
						timeoutSecs := loginTriesController.GetTimeoutLeft(r.RemoteAddr)
						if timeoutSecs != 0 {
							userLoginTry.Timeout = timeoutSecs
							userLoginTry.HasTimeout = true
						}

						//When teacher tries to login as admin, he gets same error message as if his username and password didn't match. That's the case after all.
						err := templates["login_error.gtpl"].Execute(w, userLoginTry)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							log.Print(err)
						}

						log.Print("LOGIN|Unsuccessful try to log in as admin (no permissions) from:" + username)
						return
					}
					log.Print("LOGIN|Successful admin logon from:" + r.RemoteAddr)
				}
				loginTriesController.ResetTries(r.RemoteAddr)

				log.Print("LOGIN|Logon as: " + username + " from:" + r.RemoteAddr)

				//We always create new session for users who don't have valid cookies.
				sessionID := sessionManager.GetSessionID(username)

				//Send cookie with defined expiration time and sessionID value to user.
				cookie := &http.Cookie{Name: "sessionID", Value: sessionID, Expires: time.Now().Add(cookieLifeTime * time.Second), Secure: true, HttpOnly: false}

				//Send cookie with sessionID to Client.
				http.SetCookie(w, cookie)

				//Redirect user to main.

				//TODO: is that correct http status?

				http.Redirect(w, r, "/main", http.StatusSeeOther)
			}
		}
	})
}
