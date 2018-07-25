package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"time"
)

//TODO: with every iteration we check whether there are any too old sessions. Is this overhead?

//User defines one logged user from session.
type User struct {
	username   string
	privileges string
}

//SessionManager describes basic SessionManager for netApp.
type SessionManager struct {
	sessionsToUsers          map[string]*Session
	sessionIDLength          int
	sessionLifePeriodSeconds time.Duration
}

//Session describes one session.
type Session struct {
	removeTime time.Time
	data       map[string]string
}

func (manager *SessionManager) removeOldSessions() {
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
func GetSessionManager(sessionIDLength int, lifePeriod time.Duration) *SessionManager {
	manager := &SessionManager{sessionsToUsers: make(map[string]*Session), sessionIDLength: 32, sessionLifePeriodSeconds: lifePeriod}

	return manager
}

//FindUserName returns matching username basing on provided sessionID or error if username can't be found.
func (manager *SessionManager) FindUserName(sessionID string) (string, error) {
	manager.removeOldSessions()
	value, ok := manager.sessionsToUsers[sessionID]

	if ok {
		return value.data["username"], nil
	}

	return "", errors.New("No user for such sessionID")
}

//IsLoggedIn returns true if such sessionID is stored and false otherwise.
func (manager *SessionManager) IsLoggedIn(sessionID string) bool {
	//If such sessionID exists return true, false otherwise.
	_, ok := manager.sessionsToUsers[sessionID]
	if ok {
		return true
	}

	return false
}

//GetSessionID returns unique sessionID for provided username.
func (manager *SessionManager) GetSessionID(username string) string {
	manager.removeOldSessions()
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
		sessionID = generateSessionID(manager.sessionIDLength) //TODO: make unique id
		_, ok := manager.sessionsToUsers[sessionID]
		if !ok {
			break
		}
	}

	manager.createSession(sessionID, username)
	return sessionID
}

func (manager *SessionManager) createSession(sessionID string, username string) {
	manager.removeOldSessions()
	manager.sessionsToUsers[sessionID] = &Session{time.Now().Add(manager.sessionLifePeriodSeconds), make(map[string]string)}
	manager.sessionsToUsers[sessionID].data["username"] = username
}

//RemoveSession removes session based on the provided sessionID.
func (manager *SessionManager) RemoveSession(sessionID string) {
	delete(manager.sessionsToUsers, sessionID)
}