package web

import (
	"log"
	"net/http"
	"strconv"
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
			ErrorHandler(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
			// http.Error(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryID, err := strconv.Atoi(r.FormValue("category"))
		if err != nil {
			ErrorHandler(w, "Error in CreatePost: Converting category ID", http.StatusBadRequest)
			return
		}
		log.Println("Received title:", title)
		log.Println("Received content:", content)
		log.Println("Received category ID:", categoryID)

		// Insert post into the database
		_, err = db.Exec("INSERT INTO Post (title, content, user_id, category_id, created_at) VALUES (?, ?, ?, ?)",
			title, content, userID, categoryID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			ErrorHandler(w, "errorInCreatePost", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

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
