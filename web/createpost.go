package web

import (
	"log"
	"net/http"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	// Check if user is logged in
	var userID int
	var err error

	// Verify session ID in the database
	data.LoggedIn, userID, err = VerifySession(r)
	if err != nil {
		log.Println("Error verifying session:", err)
	}
	log.Println("No session ID cookie found")

	// Render form if logged in, otherwise show error message
	if r.Method == http.MethodGet {
		RenderTemplate(w, "create-post", data)
		return
	}

	// Process post creation if POST method and logged in
	if r.Method == http.MethodPost {
		if !data.LoggedIn {
			http.Error(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		log.Println("Received title:", title)
		log.Println("Received content:", content)

		// Insert post into the database
		_, err = db.Exec("INSERT INTO Post (title, content, user_id, created_at) VALUES (?, ?, ?, ?)",
			title, content, userID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error creating post:", err)
			ErrorHandler(w, "errorInCreatePost", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// GetSessionID retrieves the session ID from the cookie
// func GetSessionID(r *http.Request) (string, error) {
// 	cookie, err := r.Cookie("session_id")
// 	if err != nil {

//		}
//		log.Println("Session ID:", cookie.Value)
//		return cookie.Value, nil
//	}
//
// VerifySession checks if the session ID exists in the database
func VerifySession(r *http.Request) (bool, int, error) {
	var userID int
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No session ID cookie found")
		return false, 0, err
	}

	err = db.QueryRow("SELECT user_id FROM Session WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		return false, 0, err // return false if session ID not found
	}
	return true, userID, nil
}
