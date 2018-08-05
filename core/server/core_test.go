package server_test

//BUG: What if server where this app is being tested has ports listed below already taken?

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/d0ku/e_register/core/server"
)

var (
	testCrt       = "../../config/certs/server.crt"
	testPrivKey   = "../../config/certs/server.key"
	testHTTPPort  = "1568"
	testHTTPSPort = "1234"
	testHost      = "localhost"

	testClientTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
)

func TestHTTPToHTTPSRedirect(t *testing.T) {
	mux := http.NewServeMux()
	serverHTTPS := server.GetTLSServer(testHTTPSPort, mux)
	go serverHTTPS.ListenAndServeTLS(testCrt, testPrivKey)
	defer serverHTTPS.Close()

	serverRedirect := server.GetRedirectServer(testHTTPSPort, testHTTPPort)
	go serverRedirect.ListenAndServe()
	defer serverRedirect.Close()

	client := &http.Client{Transport: testClientTransport}

	res, err := client.Get("http://" + testHost + ":" + testHTTPPort)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Error("Server did not send back information that it does not serve file. It means client was not redirected to HTTPS.")
	}
}

func TestHTTPToHTTPSRedirectDisabled(t *testing.T) {
	mux := http.NewServeMux()
	serverHTTPS := server.GetTLSServer(testHTTPSPort, mux)

	go serverHTTPS.ListenAndServeTLS(testCrt, testPrivKey)
	defer serverHTTPS.Close()

	client := &http.Client{Transport: testClientTransport}

	_, err := client.Get("http://" + testHost + ":" + testHTTPPort)

	if err == nil {
		t.Fatal("Client should throw an error, that it can not connect to server, because there is no redirect set.")
	}
}
