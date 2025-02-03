package web

import (
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// signUp handles both GET and POST requests for user registration
func SignUp(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, "signup", nil)
	case http.MethodPost:
		handleSignUpPost(w, r)
	default:
		ErrorHandler(w, "Invalid request method", http.StatusNotFound)
	}
}
// handleSignUpPost handles the sign up form submission
func handleSignUpPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if !isValidEmail(email) {
		handleSignUpError(w, "Invalid email address")
		return
	}

	isUnique, err := isUsernameOrEmailUnique(username, email)
	if err != nil {
		handleLoginError(w, "Error checking if username is unique", err)
		return
	}
	if !isUnique {
		handleSignUpError(w, "Username already exists")
		return
	}

	// Hash the password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		handleSignUpError(w, "error1InSignup")
		return
	}

	// Insert user into database
	err = insertUserIntoDB(username, email, hashedPassword)
	if err != nil {
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
// hashPassword hashes the user's password using bcrypt
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}
// insertUserIntoDB inserts the user's details into the database
func insertUserIntoDB(username, email, hashedPassword string) error {
	_, err := db.Exec("INSERT INTO User (username, email, password, created_at) VALUES (?, ?, ?, ?)",
		username, email, hashedPassword, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func handleSignUpError(w http.ResponseWriter, message string) {
	ErrorHandler(w, message, http.StatusNotFound)
	log.Println(message)
}
// isValidEmail checks if the email address is in a valid format
func isValidEmail(email string) bool {
	at := strings.Index(email, "@")
	if at <= 0 || at >= len(email)-1 {
		return false
	}
	dot := strings.LastIndex(email, ".")
	if dot <= at || dot >= len(email)-1 {
		return false
	}
	return true
}
// isUsernameOrEmailUnique checks if the username or email is unique in the database
func isUsernameOrEmailUnique(username, email string) (bool, error) {
    username = strings.ToLower(username)
    email = strings.ToLower(email)

    var count int
    err := db.QueryRow(`
        SELECT COUNT(*) 
        FROM User 
        WHERE username = ? OR email = ?`, username, email).Scan(&count)
    if err != nil {
        return false, err
    }
    return count == 0, nil  // Returns true if neither username nor email exists
}
