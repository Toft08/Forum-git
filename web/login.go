package web

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// login handles both GET and POST requests for user authentication
func Login(w http.ResponseWriter, r *http.Request, data *PageDetails) {
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
			ErrorHandler(w, "error1InLogin", http.StatusNotFound)
			return
		}
		// Verify submitted password matches stored hash
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			ErrorHandler(w, "error2InLogin", http.StatusNotFound)
			return
		}
		// sessionID := uuid.NewString()
		// http.SetCookie(w, &http.Cookie{
		// 	Name:     "session_id",
		// 	Value:    sessionID,
		// 	Expires:  time.Now().Add(24 * time.Hour),
		// 	HttpOnly: true,
		// 	Path:    "/",
		// })

		// Store session ID in database, associated with userID
		// _, err = db.Exec("INSERT INTO Session (id, user_id, created_at) VALUES (?, ?, ?)", sessionID, userID, time.Now().Format("2006-01-02 15:04:05"))
		// if err != nil {
		// 	http.Error(w, "Failed to create session", http.StatusInternalServerError)
		// 	return
		// }
		if err := createSession(w, userID); err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		data.LoggedIn = true

		http.Redirect(w, r, "/", http.StatusFound)

	}
}
func createSession(w http.ResponseWriter, userID int) error {

	_, err := db.Exec("DELETE FROM Session WHERE user_id = ?", userID)
	if err != nil {
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

	return err
}
