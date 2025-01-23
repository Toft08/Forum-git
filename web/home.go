package web

import (
	"net/http"
	"strconv"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, "error1InHomepage", "error", http.StatusNotFound)
		return
	}
	IsLoggedIn, _ := IsLoggedIn(r)

	message := r.URL.Query().Get("message")

	// Fetch posts from the database
	rows, err := db.Query("SELECT id, title, content FROM Post ORDER BY created_at DESC")
	if err != nil {
		// log.Println("Error fetching posts:", err)
		ErrorHandler(w, "error2InHomePage", "error", http.StatusNotFound)
		// http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build a list of posts
	var posts []map[string]string
	for rows.Next() {
		var id int
		var title, content string
		rows.Scan(&id, &title, &content)
		userID := strconv.Itoa(id)
		posts = append(posts, map[string]string{
			"ID":      userID,
			"Title":   title,
			"Content": content,
		})
	}

	// Pass posts to template
	RenderTemplate(w, "index", map[string]interface{}{
		"Message":    message,
		"Posts":      posts,
		"isLoggedIn": IsLoggedIn,
	})
}
