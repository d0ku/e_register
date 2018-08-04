package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/d0ku/e_register/core/sessions"
)

func TestLogOutHandlerWithCookieNoSuchSession(t *testing.T) {
	sessionManager = sessions.GetSessionManager(64, 120*time.Second)

	req, err := http.NewRequest("GET", "/logout", nil)

	cookie := &http.Cookie{Name: "sessionID", Value: "placeholder"}

	req.AddCookie(cookie)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	logoutHandler(rec, req)

	cookies := rec.Result().Cookies()

	var sessionCookie *http.Cookie

	for _, cookieLog := range cookies {
		if cookieLog.Name == "sessionID" {
			sessionCookie = cookieLog
			break
		}
	}

	if sessionCookie == nil {
		t.Error("Cookie removal not found.")
		return
	}

	if sessionCookie.Expires.After(time.Now()) {
		t.Error("Bad expiration for cookie removal.")
	}

	if rec.Code != http.StatusFound {
		t.Error("Incorrect status code of response.")
	}
}

func TestLogOutHandlerWithCookieAndValidSession(t *testing.T) {
	sessionManager = sessions.GetSessionManager(64, 120*time.Second)

	req, err := http.NewRequest("GET", "/logout", nil)

	sessionID := sessionManager.GetSessionID("test_data")

	cookie := &http.Cookie{Name: "sessionID", Value: sessionID}

	req.AddCookie(cookie)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	logoutHandler(rec, req)

	cookies := rec.Result().Cookies()

	var sessionCookie *http.Cookie

	for _, cookieLog := range cookies {
		if cookieLog.Name == "sessionID" {
			sessionCookie = cookieLog
			break
		}
	}

	if sessionCookie == nil {
		t.Error("Cookie removal not found.")
		return
	}

	if sessionCookie.Expires.After(time.Now()) {
		t.Error("Bad expiration for cookie removal.")
	}

	if rec.Code != http.StatusFound {
		t.Error("Incorrect status code of response.")
	}

	_, err = sessionManager.GetSession(sessionID)

	if err == nil {
		t.Error("Session was not correctly removed from sessionManager and still can be found after log out.")
	}
}

func TestLogOutHandlerWithoutCookie(t *testing.T) {
	parseAllTemplates("../../page/")
	sessionManager = sessions.GetSessionManager(64, 120*time.Second)

	req, err := http.NewRequest("GET", "/logout", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	logoutHandler(rec, req)

	if rec.Code != http.StatusFound {
		t.Error("Incorrect status code of response.")
	}
}
