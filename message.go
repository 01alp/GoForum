package main

import (
	"fmt"
	"regexp"
)

var rxEmail = regexp.MustCompile(`.+@.+\..+`)

type Message struct {
	Email    string
	Username string
	Password string
	Threads  []string
	Errors   map[string]string
}

func (msg *Message) ValidateLogin() bool {

	msg.Errors = make(map[string]string)

	user := fetchUserByEmail(database, msg.Email)
	if user == (User{}) {
		msg.Errors["Email"] = "No user with such email"
		fmt.Println("No user with such email")
		return len(msg.Errors) == 0
	}

	fmt.Println("user", user)

	// Get the expected password from database
	expectedPassword := user.Password

	if expectedPassword != msg.Password {
		// Handle wrong password
		msg.Errors["Password"] = "Wrong password"
	}

	return len(msg.Errors) == 0
}

func (msg *Message) ValidateRegistration() bool {

	msg.Errors = make(map[string]string)

	// check if username is not present in database
	if fetchUserByUsername(database, msg.Username) != (User{}) {
		msg.Errors["Username"] = "Username is present in database"
		fmt.Println("Username is present in database")
	}

	// check if email is correctly formated
	match := rxEmail.Match([]byte(msg.Email))

	if !match {
		msg.Errors["Email"] = "Please enter a valid email address"
	}

	// check if email is present in database
	if fetchUserByEmail(database, msg.Email) != (User{}) {
		msg.Errors["Email"] = "Email is present in database"
		fmt.Println("Email is present in database")
	}

	// check if username is correct length
	if len(msg.Username) < 4 || len(msg.Username) > 20 {
		msg.Errors["Username"] = "Invalid lenght of username"
		fmt.Println("Invalid lenght of Username")
	}

	// check if email is correct length
	if len(msg.Email) < 6 || len(msg.Email) > 500 {
		msg.Errors["Email"] = "Invalid lenght of email"
		fmt.Println("Invalid lenght of Email")
	}

	// check if password is correct length
	if len(msg.Password) < 4 || len(msg.Password) > 20 {
		msg.Errors["Password"] = "Invalid lenght of password"
		fmt.Println("Invalid lenght of Password")
	}
	return len(msg.Errors) == 0
}

func (msg *Message) ValidateThreads() bool {

	msg.Errors = make(map[string]string)

	// check if at least one thread is chosen when creating new post
	if len(msg.Threads) == 0 {
		msg.Errors["Threads"] = "Choose at least one category"
		// fmt.Println("Smth went wrong")
		// fmt.Println(msg.Errors["Threads"])
	}
	return len(msg.Errors) == 0
}

func (msg *Message) ValidateComment() bool {

	msg.Errors = make(map[string]string)
	
	// if len(msg.Threads) == 0 {
	// 	msg.Errors["Threads"] = "Choose at least one category"
	// 	// fmt.Println("Smth went wrong")
	// 	// fmt.Println(msg.Errors["Threads"])
	// }
	return len(msg.Errors) == 0
}

