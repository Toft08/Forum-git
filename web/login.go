package web

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// login handles both GET and POST requests for user authentication
func Login(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, "login", nil)
	case http.MethodPost:
		HandleLoginPost(w, r, data)
	default:
		ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func HandleLoginPost(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	userID, hashedPassword, err := GetUserCredentials(username)
	if err != nil {
		HandleLoginError(w, "error1InLogin", err)
		return
	}

	// Verify password
	if err := VerifyPassword(hashedPassword, password); err != nil {
		HandleLoginError(w, "error2InLogin", err)
		return
	}

	// Create session
	if err := CreateSession(w, userID); err != nil {
		ErrorHandler(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	data.LoggedIn = true
	http.Redirect(w, r, "/", http.StatusFound)
}

func GetUserCredentials(username string) (int, string, error) {
	var userID int
	var hashedPassword string

	err := db.QueryRow("SELECT id, password FROM User WHERE username = ?", username).Scan(&userID, &hashedPassword)
	if err != nil {
		return 0, "", err
	}
	return userID, hashedPassword, nil
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HandleLoginError(w http.ResponseWriter, message string, err error) {
	ErrorHandler(w, message, http.StatusNotFound)
	log.Println(message, err)
}

// CreateSession creates a new session for the user
func CreateSession(w http.ResponseWriter, userID int) error {
	_, err := db.Exec("DELETE FROM Session WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Error deleting session:", err)
		return err
	}

	sessionID := uuid.NewString()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Store session ID in database
	_, err = db.Exec("INSERT INTO Session (id, user_id, created_at) VALUES (?, ?, ?)",
		sessionID, userID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error storing session into database:", err)
		return err
	}

	return nil
}
