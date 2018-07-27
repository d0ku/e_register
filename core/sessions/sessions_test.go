package sessions_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/d0ku/database_project_go/core/sessions"
)

var manager *sessions.SessionManager

func setUpManager() {
	manager = sessions.GetSessionManager(8, 15*time.Minute)
}

func TestCreateSessionManager(t *testing.T) {
	setUpManager()

	if manager == nil {
		t.Errorf("Could not create session manager with correct parameters.")
	}

	if manager.GetSessionCount() != 0 {
		t.Errorf("Session manager should be initialized with none active sessions.")
	}
}

func TestGetSessionID(t *testing.T) {

	//TODO: this test does not make much sense at the moment.
	manager = sessions.GetSessionManager(1, 15*time.Minute)

	temp := make([]string, 4)

	for i := 0; i < 4; i++ {
		temp[i] = manager.GetSessionID("does_not_matter")
		fmt.Println(temp[i])
	}

	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			if temp[i] == temp[j] {
				t.Fatalf("Sessions IDs must be unique and they are not.")
			}
		}
	}
}
