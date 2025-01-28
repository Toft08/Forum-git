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
		handleSignUpGet(w)
	case http.MethodPost:
		handleSignUpPost(w, r)
	default:
		ErrorHandler(w, "Invalid request method", "error", http.StatusNotFound)
	}
}

func handleSignUpGet(w http.ResponseWriter) {
	RenderTemplate(w, "signup", nil)
}

func handleSignUpPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if !isValidEmail(email) {
		handleSignUpError(w, "Invalid email address")
		return
	}

	isUnique, err := isUsernameUnique(username)
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
		handleSignUpError(w, "error2InSignup")
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func insertUserIntoDB(username, email, hashedPassword string) error {
	_, err := db.Exec("INSERT INTO User (username, email, password, created_at) VALUES (?, ?, ?, ?)",
		username, email, hashedPassword, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func handleSignUpError(w http.ResponseWriter, message string) {
	ErrorHandler(w, message, "error", http.StatusNotFound)
	log.Println(message)
}

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

func isUsernameUnique(username string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM User WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
