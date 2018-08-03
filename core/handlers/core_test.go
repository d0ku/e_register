package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/d0ku/e_register/core/sessions"
)

func TestGetSessionFromRequestNoSuchCookie(t *testing.T) {
	sessionManager = sessions.GetSessionManager(64, 120*time.Second)

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	_, err = getSessionFromRequest(rec, req)

	if err == nil {
		t.Error("There is no cookie with sessionID value, so that function should error out.")
	}

	if err != ErrNoSuchSessionCookie {
		t.Error("Uncorrect error is returned.")
	}
}

func TestGetSessionFromRequestNoSuchSession(t *testing.T) {
	sessionManager = sessions.GetSessionManager(64, 120*time.Second)

	req, err := http.NewRequest("GET", "/", nil)

	cookie := &http.Cookie{Name: "sessionID", Value: "test_value"}
	req.AddCookie(cookie)

	rec := httptest.NewRecorder()

	_, err = getSessionFromRequest(rec, req)

	if err == nil {
		t.Error("There is no cookie with sessionID value, so that function should error out.")
	}

	if err != ErrNoSuchSession {
		t.Error("Uncorrect error is returned.")
	}
}

func TestGetSessionFromRequestValidSession(t *testing.T) {
	sessionManager = sessions.GetSessionManager(64, 120*time.Second)

	sessionID := sessionManager.GetSessionID("test_session")

	req, err := http.NewRequest("GET", "/", nil)

	cookie := &http.Cookie{Name: "sessionID", Value: sessionID}
	req.AddCookie(cookie)

	rec := httptest.NewRecorder()

	session, err := getSessionFromRequest(rec, req)

	if err != nil {
		t.Error("There should not be any errors, as sessionID cookie is valid and sessionManager knows about session.")
	}

	sessionOriginal, _ := sessionManager.GetSession(sessionID)

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
	sessionManager = sessions.GetSessionManager(64, 120*time.Second)
	parseAllTemplates("../../page/")

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	finalHandler := redirectToLoginPageIfUserNotLogged(http.HandlerFunc(placeHolderHandlerFunc))

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
	sessionManager = sessions.GetSessionManager(64, 120*time.Second)
	parseAllTemplates("../../page/")

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	cookie := &http.Cookie{Name: "sessionID", Value: "placeholder"}

	req.AddCookie(cookie)

	rec := httptest.NewRecorder()

	finalHandler := redirectToLoginPageIfUserNotLogged(http.HandlerFunc(placeHolderHandlerFunc))

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
	sessionManager = sessions.GetSessionManager(64, 120*time.Second)
	parseAllTemplates("../../page/")

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	sessionID := sessionManager.GetSessionID("placeholder")

	cookie := &http.Cookie{Name: "sessionID", Value: sessionID}

	req.AddCookie(cookie)

	rec := httptest.NewRecorder()

	finalHandler := redirectToLoginPageIfUserNotLogged(http.HandlerFunc(placeHolderHandlerFunc))

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
