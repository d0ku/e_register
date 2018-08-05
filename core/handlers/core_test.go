package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/d0ku/e_register/core/logging"
	"github.com/d0ku/e_register/core/sessions"
)

func setUp() *AppContext {
	//Parse all HTML templates from provided directory.
	templates := parseAllTemplates("../../page/")

	//Initialize session manager.
	sessionManager := sessions.GetSessionManager(32, 150*time.Second)

	app := &AppContext{sessionManager, templates, nil, 150 * time.Second, sessions.GetLoginTriesController()}

	return app
}

func TestGetSessionFromRequestNoSuchCookie(t *testing.T) {
	app := setUp()

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	_, err = app.getSessionFromRequest(rec, req)

	if err == nil {
		t.Error("There is no cookie with sessionID value, so that function should error out.")
	}

	if err != ErrNoSuchSessionCookie {
		t.Error("Uncorrect error is returned.")
	}
}

func TestGetSessionFromRequestNoSuchSession(t *testing.T) {
	app := setUp()
	req, err := http.NewRequest("GET", "/", nil)

	cookie := &http.Cookie{Name: "sessionID", Value: "test_value"}
	req.AddCookie(cookie)

	rec := httptest.NewRecorder()

	_, err = app.getSessionFromRequest(rec, req)

	if err == nil {
		t.Error("There is no cookie with sessionID value, so that function should error out.")
	}

	if err != ErrNoSuchSession {
		t.Error("Uncorrect error is returned.")
	}
}

func TestGetSessionFromRequestValidSession(t *testing.T) {
	app := setUp()

	sessionID := app.sessionManager.GetSessionID("test_session")

	req, err := http.NewRequest("GET", "/", nil)

	cookie := &http.Cookie{Name: "sessionID", Value: sessionID}
	req.AddCookie(cookie)

	rec := httptest.NewRecorder()

	session, err := app.getSessionFromRequest(rec, req)

	if err != nil {
		t.Error("There should not be any errors, as sessionID cookie is valid and sessionManager knows about session.")
	}

	sessionOriginal, _ := app.sessionManager.GetSession(sessionID)

	if session != sessionOriginal {
		t.Error("Session returned from request is not same that was saved in cookie.")
	}
}

func TestRedirectToLogin(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	redirectToLogin(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Error("Incorrect HTTP Status returned.")
	}
}

func placeHolderHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test_string"))
}

func TestRedirectWithErrorToLoginNoCookie(t *testing.T) {
	app := setUp()

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	finalHandler := app.checkUserLogon(http.HandlerFunc(placeHolderHandlerFunc))

	finalHandler.ServeHTTP(rec, req)

	bodyContent, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	//It means request was not redirected to not_logged page.
	if string(bodyContent) == "test_string" {
		t.Error("Request was not redirected but it should be.")
	}
}

func TestRedirectWithErrorToLoginNoSession(t *testing.T) {
	app := setUp()

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	cookie := &http.Cookie{Name: "sessionID", Value: "placeholder"}

	req.AddCookie(cookie)

	rec := httptest.NewRecorder()

	finalHandler := app.checkUserLogon(http.HandlerFunc(placeHolderHandlerFunc))

	finalHandler.ServeHTTP(rec, req)

	bodyContent, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	//It means request was not redirected to not_logged page.
	if string(bodyContent) == "test_string" {
		t.Error("Request was not redirected but it should be.")
	}
}

func TestRedirectWithErrorToLoginValidSession(t *testing.T) {
	app := setUp()

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	sessionID := app.sessionManager.GetSessionID("placeholder")

	cookie := &http.Cookie{Name: "sessionID", Value: sessionID}

	req.AddCookie(cookie)

	rec := httptest.NewRecorder()

	finalHandler := app.checkUserLogon(http.HandlerFunc(placeHolderHandlerFunc))

	finalHandler.ServeHTTP(rec, req)

	bodyContent, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	//It means request was redirected to not_logged page.
	if string(bodyContent) != "test_string" {
		t.Error("Request was redirected but it should not be.")
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("it_works"))
}

func TestCheckPermissionShouldBeGranted(t *testing.T) {
	app := setUp()

	sessionID := app.sessionManager.GetSessionID("test_data")
	session, err := app.sessionManager.GetSession(sessionID)

	if err != nil {
		t.Fatal(err)
	}

	session.Data["user_type"] = "teacher"

	finalHandler := app.checkPermission(http.HandlerFunc(testHandler), "teacher")

	req, err := http.NewRequest("GET", "/main/teacher/", nil)

	cookie := &http.Cookie{Name: "sessionID", Value: sessionID}
	req.AddCookie(cookie)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	finalHandler.ServeHTTP(rec, req)

	text, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(text) != "it_works" {
		t.Error("No permission granted, when it should be.")
	}
}

func TestCheckPermissionShouldNotBeGrantedBadUserType(t *testing.T) {
	app := setUp()

	sessionID := app.sessionManager.GetSessionID("test_data")
	session, err := app.sessionManager.GetSession(sessionID)

	if err != nil {
		t.Fatal(err)
	}

	session.Data["user_type"] = "student"

	finalHandler := app.checkPermission(http.HandlerFunc(testHandler), "teacher")

	req, err := http.NewRequest("GET", "/main/teacher/", nil)

	cookie := &http.Cookie{Name: "sessionID", Value: sessionID}
	req.AddCookie(cookie)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	finalHandler.ServeHTTP(rec, req)

	text, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(text) == "it_works" {
		t.Error("Permission granted, when it should not be.")
	}
}

func TestCheckPermissionShouldNotBeGrantedNoSession(t *testing.T) {
	app := setUp()

	finalHandler := app.checkPermission(http.HandlerFunc(testHandler), "teacher")

	req, err := http.NewRequest("GET", "/main/teacher/", nil)

	cookie := &http.Cookie{Name: "sessionID", Value: "random_session_id"}
	req.AddCookie(cookie)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	finalHandler.ServeHTTP(rec, req)

	text, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(text) == "it_works" {
		t.Error("Permission granted, when it should not be.")
	}
}

func TestCheckPermissionShouldNotBeGrantedNoCookie(t *testing.T) {
	app := setUp()

	finalHandler := app.checkPermission(http.HandlerFunc(testHandler), "teacher")

	req, err := http.NewRequest("GET", "/main/teacher/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	finalHandler.ServeHTTP(rec, req)

	text, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(text) == "it_works" {
		t.Error("Permission granted, when it should not be.")
	}
}

func TestRoutingLogin(t *testing.T) {
	mux := logging.GetMux(http.NewServeMux())
	Initialize("../../page/", 150, mux, nil)

	req, err := http.NewRequest("GET", "/main", nil)

	if err != nil {
		t.Fatal(err)
	}

	handler, pattern := mux.Handler(req)
	fmt.Println(handler)
	fmt.Println(pattern)
}
