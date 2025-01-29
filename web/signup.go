package web

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// signUp handles both GET and POST requests for user registration
func SignUp(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	if r.Method == http.MethodGet {
		RenderTemplate(w, "signup", nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Hash the password before storing in database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Attempt to insert new user into database
		//result, err := db.Exec("INSERT INTO User (username, email, password, created_at) VALUES (?, ?, ?, ?)",
		_, err = db.Exec("INSERT INTO User (username, email, password, created_at) VALUES (?, ?, ?, ?)",
			username, email, hashedPassword, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error inserting user:", err)
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Loggin the registered user in directly OR redirecting to the login page
		// userID, err := result.LastInsertId()
		// if err != nil {
		// 	log.Println("Error getting user ID:", err)
		// 	http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
		// 	return
		// }

		// if err := createSession(w, int(userID)); err != nil {
		// 	http.Error(w, "Failed to create session", http.StatusInternalServerError)
		// 	return
		// }

		// data.LoggedIn = true

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
