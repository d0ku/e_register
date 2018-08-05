package sessions_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/d0ku/e_register/core/sessions"
)

var manager *sessions.SessionManagerStruct

func setUpManager() {
	managerTemp := sessions.GetSessionManager(8, 15*time.Minute)
	manager = managerTemp.(*sessions.SessionManagerStruct)
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

func TestBasicSessionFunctionality(t *testing.T) {
	setUpManager()

	sessionId := manager.GetSessionID("testUser")

	session, _ := manager.GetSession(sessionId)
	session.Data["test_data"] = "this_is_test"

	sessionOneMoreTime, _ := manager.GetSession(sessionId)

	if sessionOneMoreTime.Data["test_data"] != "this_is_test" {
		t.Errorf("Session data is not stored properly.")
	}
}

func CheckIfSessionWillBeDeletedAfterTime(t *testing.T) {
	setUpManager()
	manager.SessionLifePeriodSeconds = 1

	id := manager.GetSessionID("test")

	time.Sleep(2)

	_, err := manager.GetSession(id)

	if err == nil {
		t.Errorf("Session should be deleted by now.")
	}
}

func CheckIfSessionWillBeDeletedAfterTimeNotCalledSession(t *testing.T) {
	setUpManager()
	manager.SessionLifePeriodSeconds = 1

	id := manager.GetSessionID("test")
	idTwo := manager.GetSessionID("testTwo")

	time.Sleep(3)

	_, err := manager.GetSession(id)

	if err == nil {
		t.Errorf("Session should be deleted by now.")
	}

	idThree := manager.GetSessionID("testThree")
	_, err = manager.GetSession(idTwo)

	if err == nil {
		t.Errorf("Session two should also be deleted by now.")
	}

	_, err = manager.GetSession(idThree)

	if err != nil {
		t.Errorf("That session should be available.")
	}
}

func TestGetSessionID(t *testing.T) {

	//TODO: this test does not make much sense at the moment.
	managerTemp := sessions.GetSessionManager(8, 15*time.Minute)
	manager = managerTemp.(*sessions.SessionManagerStruct)

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
