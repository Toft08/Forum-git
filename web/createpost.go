package web

import (
	"log"
	"net/http"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render the "Create Post" form
		renderTemplate(w, "create-post", nil)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form values
		title := r.FormValue("title")
		content := r.FormValue("content")
		// Get session ID from cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
			return
		}

		// Retrieve user ID from session
		var userID int
		err = db.QueryRow("SELECT user_id FROM Session WHERE id = ?", cookie.Value).Scan(&userID)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		// Insert post into the database
		_, err = db.Exec("INSERT INTO Post (title, content, user_id, created_at) VALUES (?, ?, ?, ?)",
			title, content, userID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error creating post:", err)
			errorHandler(w, "errorInCreatePost", "error", http.StatusNotFound)
			// http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
