package sessions_test

import (
	"testing"
	"time"

	"github.com/d0ku/database_project_go/core/sessions"
)

var testPolicy = []*sessions.TimeoutObj{
	&sessions.TimeoutObj{90, 20},
	&sessions.TimeoutObj{50, 15},
	&sessions.TimeoutObj{40, 10},
	&sessions.TimeoutObj{30, 3},
}

func TestBasicBehaviourLowestRule(t *testing.T) {
	controller := sessions.GetLoginTriesController()
	controller.SpecifyTimeoutPolicy(testPolicy)

	for i := 0; i < 29; i++ {
		controller.AddTry("test")
	}

	if controller.GetTimeoutLeft("test") != 0 {
		t.Errorf("There should be no timeout yet.")
	}

	controller.AddTry("test")

	if controller.GetTimeoutLeft("test") == 0 {
		t.Errorf("There should be timeout by now.")
	}
}

func TestBasicBehaviourResetTries(t *testing.T) {
	controller := sessions.GetLoginTriesController()
	controller.SpecifyTimeoutPolicy(testPolicy)

	for i := 0; i < 29; i++ {
		controller.AddTry("test")
	}

	if controller.GetTimeoutLeft("test") != 0 {
		t.Errorf("There should be no timeout yet.")
	}

	controller.ResetTries("test")

	controller.AddTry("test")

	if controller.GetTimeoutLeft("test") != 0 {
		t.Errorf("There should be no timeout yet (tries were reset one line above.")
	}
}

func TestBasicBehaviourAddTryWhenThereAlreadyIsTimeout(t *testing.T) {
	controller := sessions.GetLoginTriesController()
	controller.SpecifyTimeoutPolicy(testPolicy)

	for i := 0; i < 30; i++ {
		controller.AddTry("test")
	}

	timeoutLeft := controller.GetTimeoutLeft("test")

	time.Sleep(2 * time.Second)
	controller.AddTry("test")

	if controller.GetTimeoutLeft("string") > timeoutLeft {
		t.Errorf("Time should not be reset when there already is timeout.")
	}
}

func TestBasicBehaviourTimeoutEnded(t *testing.T) {
	controller := sessions.GetLoginTriesController()
	controller.SpecifyTimeoutPolicy(testPolicy)

	for i := 0; i < 30; i++ {
		controller.AddTry("test")
	}

	time.Sleep(4 * time.Second)

	timeoutLeft := controller.GetTimeoutLeft("test")

	if timeoutLeft != 0 {
		t.Errorf("Timeout should have finished by now.")
	}
}

func TestGetLoginTriesController(t *testing.T) {
	temp := sessions.GetLoginTriesController()

	if temp == nil {
		t.Errorf("GetLoginTriesController returns nil pointer")
	}
}

func TestPolicySorting(t *testing.T) {
	handler := sessions.GetLoginTriesController()

	customPolicyNotSorted := []*sessions.TimeoutObj{
		&sessions.TimeoutObj{30, 1},
		&sessions.TimeoutObj{50, 1},
		&sessions.TimeoutObj{40, 1},
		&sessions.TimeoutObj{90, 1},
	}

	handler.SpecifyTimeoutPolicy(customPolicyNotSorted)

	customPolicySorted := []*sessions.TimeoutObj{
		&sessions.TimeoutObj{90, 1},
		&sessions.TimeoutObj{50, 1},
		&sessions.TimeoutObj{40, 1},
		&sessions.TimeoutObj{30, 1},
	}

	for i, v := range handler.TimeoutPolicy {
		if *customPolicySorted[i] != *v {
			t.Errorf("Bad sorting in LoginTries Controller (policy has to be sorted so bigger values are checked first")
		}
	}
}
