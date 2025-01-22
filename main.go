package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	db = initDB()
	defer db.Close()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/login", login)

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// renderTemplate handles the rendering of HTML templates with provided data
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// homePage renders the index page
func homePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

// signUp handles both GET and POST requests for user registration
func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "signup", nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Hash the password before storing in database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Attempt to insert new user into database
		_, err = db.Exec("INSERT INTO User (username, email, password, created_at) VALUES (?, ?, ?, ?)", username, email, hashedPassword, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// login handles both GET and POST requests for user authentication
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "login", nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Query database for user's hashed password using their email
		var hashedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Verify submitted password matches stored hash
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
