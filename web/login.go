package web

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login handles both GET and POST requests for user authentication
func Login(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	data.ValidationError = ""
	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, "login", data)
	case http.MethodPost:
		HandleLoginPost(w, r, data)
	default:
		ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// HandleLoginPost handles the user login form submission
func HandleLoginPost(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	userID, hashedPassword, err := getUserCredentials(username)
	if err != nil {
		data.ValidationError = "Invalid username"
		RenderTemplate(w, "login", data)
		return
	}

	// Verify password
	if err := verifyPassword(hashedPassword, password); err != nil {
		data.ValidationError = "Invalid password"
		RenderTemplate(w, "login", data)
		return
	}

	// Create session
	if err := createSession(w, userID); err != nil {
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.LoggedIn = true
	http.Redirect(w, r, "/", http.StatusFound)
}

// getUserCredentials retrieves the user's ID and hashed password from the database
func getUserCredentials(username string) (int, string, error) {
	var userID int
	var hashedPassword string

	err := db.QueryRow("SELECT id, password FROM User WHERE username = ?", username).Scan(&userID, &hashedPassword)
	if err != nil {
		return 0, "", err
	}
	return userID, hashedPassword, nil
}

// verifyPassword compares the hashed password with the password provided by the user
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// createSession creates a new session for the user and stores it in the database
func createSession(w http.ResponseWriter, userID int) error {
	// First check for and delete any existing sessions for this user
	_, err := db.Exec("DELETE FROM Session WHERE user_id = ?", userID)
	if err != nil {
		return err
	}
	sessionID := uuid.NewString()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: true,
		Path:     "/",
	})

	// Store session ID in database
	_, err = db.Exec("INSERT INTO Session (id, user_id, created_at) VALUES (?, ?, ?)",
		sessionID, userID, time.Now().Format("2006-01-02 15:04:05"))

	return err
}
