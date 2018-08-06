package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"time"
)

//TODO: race condition!
//Old sessions are removed automatically only when there is a request for new session.
//Of course session expiration time is validated when session is queried from manager.

//User defines one logged user from session.
type User struct {
	username   string
	privileges string
}

//SessionManager is interface which defines capabilities of session manager.
type SessionManager interface {
	GetSessionCount() int
	GetSession(string) (*Session, error)
	GetSessionID(string) string
	RemoveSession(string)
}

//SessionManagerStruct describes basic SessionManagerStruct for netApp.
type SessionManagerStruct struct {
	sessionsToUsers          map[string]*Session
	SessionIDLength          int
	SessionLifePeriodSeconds time.Duration
}

//Session describes one session.
type Session struct {
	removeTime time.Time
	Data       map[string]string
}

//GetSessionCount returns current amount of sessions.
func (manager *SessionManagerStruct) GetSessionCount() int {
	return len(manager.sessionsToUsers)
}

func (manager *SessionManagerStruct) removeOldSessions() {
	//remove all sessions whose lifespan ended.
	timeNow := time.Now()
	for index, session := range manager.sessionsToUsers {
		if timeNow.After(session.removeTime) {
			manager.RemoveSession(index)
			log.Print("Session timed out: " + index)
		}
	}
}

//GetSessionManager returns Session Manager with properly set up attributes.
func GetSessionManager(sessionIDLength int, lifePeriod time.Duration) SessionManager {
	manager := &SessionManagerStruct{sessionsToUsers: make(map[string]*Session), SessionIDLength: sessionIDLength, SessionLifePeriodSeconds: lifePeriod}

	sessionManager := SessionManager(manager)

	return sessionManager
}

//GetSessionID returns unique sessionID for provided username.
func (manager *SessionManagerStruct) GetSessionID(username string) string {
	generateSessionID := func(length int) string {
		bytes := make([]byte, length)

		_, err := io.ReadFull(rand.Reader, bytes)

		if err != nil {
			log.Fatal(err)
		}

		return base64.URLEncoding.EncodeToString(bytes)
	}

	var sessionID string
	for {
		//Try to generate unique session id until there is no session with such id. possible BUG: check for race condition
		sessionID = generateSessionID(manager.SessionIDLength)
		_, ok := manager.sessionsToUsers[sessionID]
		if !ok {
			break
		}
	}

	manager.createSession(sessionID, username)
	log.Print("New session created for:" + username)
	return sessionID
}

//GetSession returns session coupled with provided sessionID.
func (manager *SessionManagerStruct) GetSession(sessionID string) (*Session, error) {
	session, ok := manager.sessionsToUsers[sessionID]

	if !ok {
		return nil, errors.New("There is no session with such sessionID")
	}

	//check if session is still valid.

	if session.removeTime.Before(time.Now()) {
		delete(manager.sessionsToUsers, sessionID)
		return nil, errors.New("Session is outdated and has to be deleted")
	}

	return session, nil
}

func (manager *SessionManagerStruct) createSession(sessionID string, username string) {
	manager.removeOldSessions()
	manager.sessionsToUsers[sessionID] = &Session{time.Now().Add(manager.SessionLifePeriodSeconds), make(map[string]string)}
	manager.sessionsToUsers[sessionID].Data["username"] = username
}

//RemoveSession removes session based on the provided sessionID.
func (manager *SessionManagerStruct) RemoveSession(sessionID string) {
	session, err := manager.GetSession(sessionID)
	if err != nil {
		log.Print("Try to remove session which already does not exist.")
	} else {
		log.Print("Removed session of user:" + session.Data["user_name"])
		delete(manager.sessionsToUsers, sessionID)
	}
}
