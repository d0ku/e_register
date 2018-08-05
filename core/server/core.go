package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/d0ku/e_register/core/logging"
)

func redirectToHTTPS(h http.Handler, ports ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		host := r.Host
		if i := strings.Index(host, ":"); i != -1 {
			host = host[:i]
		}
		redirectAddress := "https://" + host + ":" + ports[0] + r.RequestURI

		log.Print("HTTPS REDIRECT|Redirected " + r.RemoteAddr + " to " + redirectAddress)
		http.Redirect(w, r, redirectAddress, http.StatusMovedPermanently)
	})
}

func placeHolderHandler(w http.ResponseWriter, r *http.Request) {
}

//GetTLSServer returns server with set up Port.
func GetTLSServer(HTTPSPort string, mux http.Handler) *http.Server {
	server := &http.Server{
		Addr:    ":" + HTTPSPort,
		Handler: mux,
	}

	return server
}

//GetRedirectServer returns server used to redirect http request to https.
func GetRedirectServer(HTTPSPort string, HTTPPort string) *http.Server {
	redirectHandler := logging.LogRequests(redirectToHTTPS(http.HandlerFunc(placeHolderHandler), HTTPSPort))

	server := &http.Server{
		Addr:    ":" + HTTPPort,
		Handler: redirectHandler,
	}

	return server
}
