package web

import (
	"log"
	"net/http"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	var isLoggedIn bool
	cookie, err := r.Cookie("session_id")
	if err == nil {
		log.Println("Session ID:", cookie.Value)
		// Check if session ID exists in the database
		var sessionID string
		err = db.QueryRow("SELECT id FROM Session WHERE id = ?", cookie.Value).Scan(&sessionID)
		if err == nil && sessionID != "" {
			isLoggedIn = true
		}
	} else {
		log.Println("No session ID cookie found")
	}

	// Render form if logged in, otherwise show error message
	if r.Method == http.MethodGet {
		renderTemplate(w, "create-post", map[string]interface{}{
			"IsLoggedIn": isLoggedIn,
		})
		return
	}

	// Process post creation if POST method and logged in
	if r.Method == http.MethodPost {
		if !isLoggedIn {
			http.Error(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		log.Println("Received title:", title)
		log.Println("Received content:", content)

		// Get user ID from session
		var userID int
		err = db.QueryRow("SELECT user_id FROM Session WHERE id = ?", cookie.Value).Scan(&userID)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}
		log.Println("User ID for sessionn:", userID)

		// Insert post into the database
		_, err = db.Exec("INSERT INTO Post (title, content, user_id, created_at) VALUES (?, ?, ?, ?)",
			title, content, userID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error creating post:", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
