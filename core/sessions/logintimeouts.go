package sessions

import (
	"log"
	"math"
	"sort"
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
	tries         map[string]int
	timeoutEnd    map[string]time.Time
	TimeoutPolicy []*TimeoutObj
}

//TimeoutObj represents one rule from timeoutPolicy.
type TimeoutObj struct {
	HowManytries     int
	HowLongInSeconds int
}

//GetLoginTriesController returns pointer to initalized controller.
func GetLoginTriesController() *LoginTriesController {
	defaultTimeoutPolicy := []*TimeoutObj{
		&TimeoutObj{100, 30},
		&TimeoutObj{50, 20},
		&TimeoutObj{10, 15},
		&TimeoutObj{5, 10},
	}
	return &LoginTriesController{make(map[string]int), make(map[string]time.Time), defaultTimeoutPolicy}
}

//SpecifyTimeoutPolicy makes it very easy to use your own rules for timeout.
//User has to provide slice with rules defined as TimeoutObj objects.
//Function sorts that slice and uses it accordingly.
func (controller *LoginTriesController) SpecifyTimeoutPolicy(policy []*TimeoutObj) {
	sort.Slice(policy, func(i, j int) bool { return policy[i].HowManytries > policy[j].HowManytries })
	controller.TimeoutPolicy = policy
}

func (controller *LoginTriesController) setTimeout(origin string) {
	var addTime time.Duration
	tries := controller.tries[origin]

	for _, value := range controller.TimeoutPolicy {
		if tries >= value.HowManytries {
			addTime = time.Second * time.Duration(value.HowLongInSeconds)
			break
		}
	}

	log.Print(addTime.String() + " seconds of timeout for:" + origin)

	controller.timeoutEnd[origin] = time.Now().Add(addTime)
}

//AddTry add one try to user, if tries count exceeds specified value, it calls setTimeout function which set timeouts according to specified policy.
func (controller *LoginTriesController) AddTry(origin string) {
	_, ok := controller.timeoutEnd[origin]
	if ok {
		//Don't do anything if there is a timeout already. That request won't be supported anyway.
		return
	}
	_, ok = controller.tries[origin]

	if !ok {
		controller.tries[origin] = 0
	}

	controller.tries[origin]++

	if controller.tries[origin] >= controller.TimeoutPolicy[len(controller.TimeoutPolicy)-1].HowManytries {
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

	currentTime := time.Now()
	if currentTime.After(timeout) {
		//Timeout already passed, can be deleted.
		delete(controller.timeoutEnd, origin)
		return 0
	}

	//Return time left.
	timeLeft := int(math.Ceil(timeout.Sub(currentTime).Seconds()))
	return timeLeft
}
