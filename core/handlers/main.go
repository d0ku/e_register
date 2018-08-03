package handlers

//All requests that are called in this file should already pass through redirectToLogin handler decorator, so we can be sure that they have cookie with correct value set up.

import (
	"log"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, err := getSessionFromRequest(w, r)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = templates["index.gtpl"].Execute(w, session.Data["username"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}
