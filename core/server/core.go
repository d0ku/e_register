package server

import (
	"fmt"
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

//RunTLS starts initialized server on specified port with TLS.
func RunTLS(HTTPSport string, HTTPort string, redirectHTTPtoHTTPS bool, hostname string, serverCert string, serverKey string) {

	//Redirect HTTP trafic to HTTPS port with changed protocol if such option was specified.
	if redirectHTTPtoHTTPS {
		go func() {
			err := http.ListenAndServe(":"+HTTPort, logging.LogRequests(redirectToHTTPS(http.HandlerFunc(placeHolderHandler), HTTPSport)))
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	fmt.Println()
	fmt.Println("Listen to me at: https://" + hostname + ":" + HTTPSport)

	err := http.ListenAndServeTLS(":"+HTTPSport, serverCert, serverKey, nil)

	//Something went wrong with starting HTTPS server.
	if err != nil {
		panic(err)
	}
}
