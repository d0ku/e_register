package logging

import (
	"log"
	"net/http"
	"strconv"
)

//ResponseWriter is abstraction over http.ResponseWriter which logs response status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

//GetResponseWriter returns one object capable of catching http response status code and logging it.
func getResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

//WriteHeader is basic abstraction over http.ResponseWriter.WriteHeader which logs response status code.
func (writer *responseWriter) writeHeader(code int) {
	writer.statusCode = code
	writer.ResponseWriter.WriteHeader(code)
}

//LogRequests is decorator for http.Handler functions which logs all incoming requests and outcoming responses.
func LogRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Print("REQ|" + r.RemoteAddr + "|" + r.Method + "|" + r.URL.String())

		responseWriter := getResponseWriter(w)
		handler.ServeHTTP(responseWriter, r)

		status := responseWriter.statusCode
		log.Print("RES|" + strconv.Itoa(status) + "|" + http.StatusText(status))
	})
}
