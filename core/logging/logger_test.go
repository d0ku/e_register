package logging_test

//TODO: that test does not test anything at the moment.
import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/d0ku/e_register/core/logging"
)

func TestCatchingHandlerResponse(t *testing.T) {

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	finalHandler := logging.LogRequests(http.HandlerFunc(testHandler))

	finalHandler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusTooManyRequests {

	}

}
