package core

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
)

//TODO: when session is created, we have to make sure that id will be unique.

//User defines one logged user from session.
type User struct {
	username   string
	privileges string
}

//SessionManager describes basic SessionManager for netApp.
type SessionManager struct {
	sessionsToUsers map[string]string
	usersToSessions map[string]string
	sessionIDLength int
	//If no dbConnection is provided, there won't be any operations done in SQL.
	dbConnection *DBHandler
}

//GetSessionManager returns Session Manager with properly set up attributes.
func GetSessionManager(sessionIDLength int, dbConnection *DBHandler) *SessionManager {
	manager := &SessionManager{sessionsToUsers: make(map[string]string), usersToSessions: make(map[string]string), sessionIDLength: 32, dbConnection: dbConnection}

	return manager
}

//FindSessionID returns matching sessionID basing on provided username or error if sessionID can't be found.
func (manager *SessionManager) FindSessionID(username string) (string, error) {
	value, ok := manager.usersToSessions[username]

	if ok {
		return value, nil
	}

	return "", errors.New("No sessionID for such username")
}

//FindUserName returns matching username basing on provided sessionID or error if username can't be found.
func (manager *SessionManager) FindUserName(sessionID string) (string, error) {
	value, ok := manager.sessionsToUsers[sessionID]

	if ok {
		return value, nil
	}

	return "", errors.New("No user for such sessionID")
}

//GetSessionID returns already existing sessionID for username, or creates new and returns if needed.
func (manager *SessionManager) GetSessionID(username string) string {

	generateSessionID := func(length int) string {
		bytes := make([]byte, length)

		_, err := io.ReadFull(rand.Reader, bytes)

		if err != nil {
			log.Fatal(err)
		}

		return base64.URLEncoding.EncodeToString(bytes)
	}

	sessionID, err := manager.FindSessionID(username)

	if err == nil {
		return sessionID
	}

	sessionID = generateSessionID(manager.sessionIDLength) //TODO: make unique id
	manager.createSession(sessionID, username)
	return sessionID

}

func (manager *SessionManager) createSession(sessionID string, username string) {
	manager.usersToSessions[username] = sessionID
	manager.sessionsToUsers[sessionID] = username

	if manager.dbConnection != nil {
		manager.dbConnection.AddSession(sessionID, username)
	}
}

//RemoveSession removes session based on the provided sessionID.
func (manager *SessionManager) RemoveSession(sessionID string) {
	delete(manager.sessionsToUsers, manager.sessionsToUsers[manager.usersToSessions[sessionID]])
	delete(manager.usersToSessions, manager.usersToSessions[sessionID])
	if manager.dbConnection != nil {
		manager.dbConnection.DeleteSession(sessionID)
	}
}

//ReadSessionsFromDatabase tries to find session_id->username pairs in provided database in table SessionID.
func (manager *SessionManager) ReadSessionsFromDatabase() error {
	if manager.dbConnection == nil {
		return errors.New("No database handler provided")
	}

	manager.sessionsToUsers = make(map[string]string)
	rows, err := manager.dbConnection.Query("SELECT * FROM SessionID;")
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	var sessionID string
	var username string
	fmt.Println("Restoring sessions:")
	for rows.Next() {
		rows.Scan(&sessionID, &username)
		manager.sessionsToUsers[sessionID] = username
		manager.usersToSessions[username] = sessionID
		fmt.Println("---> " + username)
	}
	fmt.Println("Sessions restored.")
	return nil
}
