package core

//TODO: sessions are now stored in database, they should be removed after certain amount of time.
import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	//That's recommended way to import sql driver.
	_ "github.com/lib/pq"
)

var dbConnection *sql.DB

var sessionManager *SessionManager
var templates map[string]*template.Template

func parseAllTemplates() {
	pageFolder := "./page/"
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
	//That implementation is not stable and dangerous, but fun so I will keep it for now.

	//If address is "/" point to index, else try to find adequate file in page folder.

	fmt.Println(request.Method)
	if request.Method == "GET" {
		//		t, _ := template.ParseFiles("page/index.gtpl")
		err := templates["index.gtpl"].Execute(response, nil)
		if err != nil {
			log.Fatal(err)
		}

		//Set correct Content-Type in response basing on request.

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

		username, err := sessionManager.FindUserName(cookie.Value)

		if err != nil {
			log.Print("Client side thinks it is logged in, but it is not.")

			//TODO: display something about that he has to relogin (probably redirect should be enough).
			cookie.Expires = time.Unix(0, 0)
			http.SetCookie(response, cookie)
			http.Redirect(response, request, "/login", http.StatusSeeOther)
			//user thinks he is logged in, but it is not true.
			return

		}

		err = templates["login_personal.gtpl"].Execute(response, username)
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
		fmt.Println(username)
		fmt.Println(password)
		if len(username) < 1 || len(password) < 8 {
			//just display that no user found or that it is too empty?
			fmt.Println("Too short password!")
			return
		}

		var usernameBase sql.NullString
		err := dbConnection.QueryRow("SELECT check_login_data('" + username + "','" + password + "')").Scan(&usernameBase)
		switch {
		case err == sql.ErrNoRows || !usernameBase.Valid:
			log.Printf("No user with matching username and password.")
			return
		case err != nil:
			log.Fatal(err)
			return
		default:
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

			http.Redirect(response, request, "/login", http.StatusSeeOther)
		}

		//If user and password don't match in database return error,
		//else create session and proceed
	}
}

func deleteHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte(request.URL.String()))
}

//Run starts the server at default port and prints info in console how to connect to it.
func Run() {
	//Initialize parsed templates.
	parseAllTemplates()

	//Initialize DB connection.
	//TODO: change user to something more secure (non-root).
	var err error
	//TODO: connect to postgres by SSL (sslmode=verify-full)
	connStr := "user=d0ku dbname=database_project_go sslmode=disable"
	dbConnection, err = sql.Open("postgres", connStr)

	sessionManager = GetSessionManager(32, dbConnection)
	sessionManager.ReadSessionsFromDatabase()

	//Could not initialize connection.
	if err != nil {
		panic(err)
	}

	defer dbConnection.Close()

	var port = "1234"

	fmt.Println()
	fmt.Println("Listen to me at: https://localhost:" + port)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/main", mainHandler)
	http.Handle("/", http.FileServer(http.Dir("./page/")))
	err = http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil)

	//Something went wrong with starting HTTPS server.
	if err != nil {
		panic(err)
	}
}
