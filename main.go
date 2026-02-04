package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 1. The Big Book (Where we keep users)
var userDB = map[string]string{}

// This is what a user looks like when they send data
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// 2. Tell the Guard what to do when people visit specific doors
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)

	fmt.Println("Server started on :8080 (The Clubhouse is open!)")
	// This starts the server
	http.ListenAndServe(":8080", nil)
}

// 3. The Sign Up Rule
func signupHandler(w http.ResponseWriter, r *http.Request) {
	// We only allow POST (sending data), not GET (just looking)
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "I don't understand!", http.StatusBadRequest)
		return
	}

	userDB[newUser.Username] = newUser.Password
	fmt.Fprintf(w, "Welcome to the club, %s!", newUser.Username)
	fmt.Println("New user signed up:", newUser.Username) // Log to terminal
}

// 4. The Login Rule
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var tryingUser User
	err := json.NewDecoder(r.Body).Decode(&tryingUser)
	if err != nil {
		http.Error(w, "I don't understand!", http.StatusBadRequest)
		return
	}

	savedPassword, exists := userDB[tryingUser.Username]

	if !exists || savedPassword != tryingUser.Password {
		http.Error(w, "Wrong secret whisper! Go away!", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Open Sesame! You are logged in, %s!", tryingUser.Username)
}