package sessions

import (
	"log"
	"math"
	"time"
)

// TODO: Is this memory safe to never purge addresses which didn't login successfully?

//UserLoginTry is basic type that describes login try from an user.
type UserLoginTry struct {
	UserType   string
	UserName   string
	Timeout    int
	HasTimeout bool
}

//LoginTriesController enables application to easily control how many times user tried to log in and give him specific timeouts.
type LoginTriesController struct {
	tries      map[string]int
	timeoutEnd map[string]time.Time
}

//GetLoginTriesController returns pointer to initalized controller.
func GetLoginTriesController() *LoginTriesController {
	return &LoginTriesController{make(map[string]int), make(map[string]time.Time)}
}

func (controller *LoginTriesController) setTimeout(origin string) {
	var addTime time.Duration
	tries := controller.tries[origin]

	//TODO: define better policy for timeouts.
	if tries >= 100 {
		addTime = time.Second * 30
	} else if tries >= 50 {
		addTime = time.Second * 20
	} else if tries >= 10 {
		addTime = time.Second * 15
	} else if tries >= 5 {
		addTime = time.Second * 10
	}

	log.Print(addTime.String() + " seconds of timeout for:" + origin)

	controller.timeoutEnd[origin] = time.Now().Add(addTime)
}

//AddTry add one try to user, if tries count exceeds specified value, it calls setTimeout function which set timeouts according to specified policy.
func (controller *LoginTriesController) AddTry(origin string) {
	_, ok := controller.timeoutEnd[origin]
	if ok {
		//Don't do anything if there is a timeout already.
		return
	}
	_, ok = controller.tries[origin]

	if !ok {
		controller.tries[origin] = 0
	}

	controller.tries[origin]++

	if controller.tries[origin] >= 5 {
		controller.setTimeout(origin)
	}
}

//ResetTries resets try counter.
func (controller *LoginTriesController) ResetTries(origin string) {
	delete(controller.tries, origin)
}

//GetTimeoutLeft returns 0 if there is no timeout left, or int (seconds) representing how long user has to wait.
//If timeout left is equal to 0, user can try to log in.
//Else function should return timeout left (in seconds).
func (controller *LoginTriesController) GetTimeoutLeft(origin string) int {
	timeout, ok := controller.timeoutEnd[origin]

	if !ok {
		//There is no timeout set.
		return 0
	}

	if time.Now().After(timeout) {
		//Timeout already passed, can be deleted.
		delete(controller.timeoutEnd, origin)
		return 0
	}

	//Return time left.
	timeLeft := int(math.Round(timeout.Sub(time.Now()).Seconds()))
	if timeLeft == 0 {
		delete(controller.timeoutEnd, origin)
		return 0
	}

	return timeLeft
}
