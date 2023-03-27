package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func setSessionToken(w http.ResponseWriter, creds Credentials) {
	// Create a new random session token
	uuid, _ := uuid.NewV4()
	sessionToken := (uuid).String()
	expiresAt := time.Now().Add(15 * 60 * time.Second)
	// Set the token in the session map, along with the session information
	sessions[sessionToken] = session{
		user:   fetchUserByEmail(database, creds.Email),
		expiry: expiresAt,
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
}

func setLastPage(w http.ResponseWriter, url string) {
	expiresAt := time.Now().Add(60 * 60 * time.Second)
	http.SetCookie(w, &http.Cookie{
		Name:    "last_page",
		Value:   url,
		Expires: expiresAt,
	})
}

func welcome(w http.ResponseWriter, r *http.Request) Data {

	output := Data{LoggedIn: false, User: User{}, Threads: fetchAllThreads(database)}

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			fmt.Println("Unauthorized")
			return output
		}
		// For any other type of error, return a bad request status
		fmt.Println("Bad Request")
		return output
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		fmt.Println("Unauthorized")
		return output
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		fmt.Println("Unauthorized")
		return output
	}
	// If the session is valid, return the welcome message to the user
	output = Data{LoggedIn: true, User: userSession.user}
	fmt.Printf("\nWelcome %s!\n", userSession.user.Username)
	return output
}

func refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			fmt.Println("Unauthorized")
		}
		// For any other type of error, return a bad request status
		fmt.Println("Bad Request")
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		fmt.Println("Unauthorized")
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		fmt.Println("Unauthorized")
	}
	// If the previous session is valid, create a new session token for the current user
	uuid, _ := uuid.NewV4()
	newSessionToken := (uuid).String()
	expiresAt := time.Now().Add(120 * time.Second)
	// Set the token in the session map, along with the user whom it represents
	sessions[newSessionToken] = session{
		user:   userSession.user,
		expiry: expiresAt,
	}
	// Delete the older session token
	delete(sessions, sessionToken)
	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(15 * 60 * time.Second),
	})
}

// add refresh func before every action