package web

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// login handles both GET and POST requests for user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RenderTemplate(w, "login", nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Query database for user's hashed password using their email
		var hashedPassword string
		var userID int
		err := db.QueryRow("SELECT id, password FROM User WHERE username = ?", username).Scan(&userID, &hashedPassword)
		if err != nil {
			ErrorHandler(w, "error1InLogin", "error", http.StatusNotFound)
			// http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Verify submitted password matches stored hash
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			ErrorHandler(w, "error2InLogin", "error", http.StatusNotFound)
			// http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/?message=Login successful!", http.StatusFound)

	}
}
