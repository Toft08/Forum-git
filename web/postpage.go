package web

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PostHandler handles requests to view a specific post
func PostHandler(w http.ResponseWriter, r *http.Request, data *PageDetails) {

	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		log.Println("Error converting postID to int:", err)
		ErrorHandler(w, "Page Not Found", http.StatusNotFound)
	}

	// Check the cookie and get userID
	var userID int
	data.LoggedIn, userID = VerifySession(r)
	data.Posts = nil

	if r.Method == http.MethodPost {
		data.LoggedIn, _ = VerifySession(r)

		if !data.LoggedIn {
			ErrorHandler(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
			return
		}

		content := r.FormValue("comment")

		// Insert post into the database
		_, err := db.Exec("INSERT INTO Comment (post_id, content, user_id, created_at) VALUES (?, ?, ?, ?)",
			postID, content, userID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error creating post:", err)
			ErrorHandler(w, "errorInCreatePost", http.StatusNotFound)
			return
		}
	} else if r.Method != http.MethodGet {
		ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	post, err := GetPostDetails(postID)
	if err != nil {
		log.Println("Error fetching post details:", err)
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Posts = append(data.Posts, *post)

	RenderTemplate(w, "post", data)

}
