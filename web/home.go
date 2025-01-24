package web

import (
	"log"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, "error1InHomepage", "error", http.StatusNotFound)
		return
	}

	data := PageDetails{}

	IsLoggedIn, _ := IsLoggedIn(r)

	data.LoggedIn = IsLoggedIn

	//message := r.URL.Query().Get("message")

	// Fetch posts from the database
	rows, err := db.Query("SELECT id FROM Post ORDER BY created_at DESC")
	if err != nil {
		// log.Println("Error fetching posts:", err)
		ErrorHandler(w, "error2InHomePage", "error", http.StatusNotFound)
		// http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build a list of posts
	posts := []PostDetails{}

	for rows.Next() {
		log.Println("Post id rows")
		var id int
		rows.Scan(&id)
		post, err := getPostDetails(id)

		if err != nil {
			ErrorHandler(w, "Internal Server Error", "error", http.StatusInternalServerError)
		}
		posts = append(posts, *post)

	}
	data.Posts = posts

	// Pass posts to template
	RenderTemplate(w, "index", data)
}
