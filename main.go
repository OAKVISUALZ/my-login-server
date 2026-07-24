package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // This is the driver for Postgres
)

// We use a global variable to hold the database connection
var db *sql.DB

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// 1. Get the Database URL from the environment
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// 2. Connect to the Database
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Create the table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}

	// 4. Register the Doors (Routes)
	http.HandleFunc("/", homeHandler)         // The Front Door (NEW!)
	http.HandleFunc("/signup", signupHandler) // The Signup Door
	http.HandleFunc("/login", loginHandler)   // The Login Door

	// 5. Port Configuration for Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started on :%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

// NEW: The Home Page Handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Clubhouse API! 🏠\nUse /signup to join.\nUse /login to enter."))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write to the Database (INSERT)
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, u.Password)
	if err != nil {
		http.Error(w, "Username likely taken", http.StatusConflict)
		return
	}

	fmt.Fprintf(w, "User %s created!", u.Username)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read from the Database (SELECT)
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=$1", u.Username).Scan(&storedPassword)

	if err != nil || storedPassword != u.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Welcome back, %s!", u.Username)
}
