package handlers

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/d0ku/e_register/core/databasehandling"
)

func TestLogOutHandlerWithCookieNoSuchSession(t *testing.T) {
	app := setUp()

	req, err := http.NewRequest("GET", "/logout", nil)

	cookie := &http.Cookie{Name: "sessionID", Value: "placeholder"}
	req.AddCookie(cookie)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	app.logoutHandler(rec, req)

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
	app := setUp()

	req, err := http.NewRequest("GET", "/logout", nil)

	sessionID := app.sessionManager.GetSessionID("test_data")

	cookie := &http.Cookie{Name: "sessionID", Value: sessionID}

	req.AddCookie(cookie)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	app.logoutHandler(rec, req)

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

	_, err = app.sessionManager.GetSession(sessionID)

	if err == nil {
		t.Error("Session was not correctly removed from sessionManager and still can be found after log out.")
	}
}

func TestLogOutHandlerWithoutCookie(t *testing.T) {
	app := setUp()

	req, err := http.NewRequest("GET", "/logout", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	app.logoutHandler(rec, req)

	if rec.Code != http.StatusFound {
		t.Error("Incorrect status code of response.")
	}
}

type dbMock struct {
}

func (dbMock) CheckIfTeacherIsSchoolAdmin(id int) int {
	return -1
}

func (dbMock) CheckUserLogin(username string, password string, userType string) *databasehandling.UserLoginData {
	if username == "test_teacher" && password == "teacher_password" && userType == "teacher" {

		return &databasehandling.UserLoginData{true, "teacher", 1, false}
	}
	return &databasehandling.UserLoginData{false, "", 0, false}
}

func (dbMock) GetSchoolsDetailsWhereTeacherTeaches(id string) ([]databasehandling.School, error) {
	return nil, nil
}

func getClientWithTurnedOffCertificateMatching() *http.Client {
	//BUG: When i PostForm to httptest.NewTLSServer with http.PostForm i get certificate denial. It should be changed in future.
	tn := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &http.Client{
		Transport: tn,
		Jar:       jar,
	}
}

func TestLogInHandlerAsTeacherMockedDBShouldNotLogIn(t *testing.T) {
	app := setUp()

	app.DbHandler = databasehandling.DBHandler(&dbMock{})

	testServer := httptest.NewTLSServer(http.HandlerFunc(app.loginHandler))
	defer testServer.Close()

	testReq, err := http.NewRequest("GET", "/main/teacher/1", nil)

	if err != nil {
		t.Fatal(err)
	}

	loginData := url.Values{
		"username": []string{"teacher"},
		"password": []string{"teacher_password"},
		"userType": []string{"teacher"}}

	resp, err := getClientWithTurnedOffCertificateMatching().PostForm(testServer.URL, loginData)

	if err != nil {
		t.Fatal(err)
	}

	for _, cookie := range resp.Cookies() {
		testReq.AddCookie(cookie)
	}

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	loggedTeacherHandler := app.checkPermission(http.HandlerFunc(testHandler), "teacher")

	loggedTeacherHandler.ServeHTTP(rec, testReq)

	text, err := ioutil.ReadAll(rec.Body)

	if err != nil {
		t.Fatal(err)
	}

	log.Print(string(text))

	if string(text) == "it_works" {
		t.Error("Permission granted when it should not be.")
	}
}

func TestLogInHandlerAsTeacherMockedDBShouldLogIn(t *testing.T) {
	//TODO: skipped because of strange cookie behaviour
	t.SkipNow()
	app := setUp()

	app.DbHandler = databasehandling.DBHandler(&dbMock{})

	testServer := httptest.NewTLSServer(http.HandlerFunc(app.loginHandler))
	defer testServer.Close()

	testReq, err := http.NewRequest("GET", "/main/teacher", nil)

	if err != nil {
		t.Fatal(err)
	}

	loginData := url.Values{
		"username": []string{"test_teacher"},
		"password": []string{"teacher_password"},
		"userType": []string{"teacher"}}

	resp, err := getClientWithTurnedOffCertificateMatching().PostForm(testServer.URL, loginData)

	if err != nil {
		t.Fatal(err)
	}

	var val []byte
	_, _ = resp.Body.Read(val)
	resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		log.Print(cookie)
		testReq.AddCookie(cookie)
	}

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	loggedTeacherHandler := app.checkPermission(http.HandlerFunc(testHandler), "teacher")

	loggedTeacherHandler.ServeHTTP(rec, testReq)

	text, err := ioutil.ReadAll(rec.Body)

	if err != nil {
		t.Fatal(err)
	}

	log.Print(string(text))

	if string(text) != "it_works" {
		t.Error("Permission not granted when it should be.")
	}
}
