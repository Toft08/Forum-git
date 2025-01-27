package web

import (
	"database/sql"
	"log"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	if r.URL.Path != "/" {
		ErrorHandler(w, "error1InHomepage", "error", http.StatusNotFound)
		return
	}

	var userID int
	var rows *sql.Rows
	var err error
	data.LoggedIn, userID = IsLoggedIn(r)

	//message := r.URL.Query().Get("message")
	if data.LoggedIn && r.Method == http.MethodPost {
		// Fetch posts from the database for a specific user
		rows, err = db.Query("SELECT id FROM Post WHERE user_id = ? ORDER BY created_at DESC", userID)
		if err != nil {
			log.Println("Error fetching users own posts:", err)
			ErrorHandler(w, "errorFetchingPosts", "error", http.StatusNotFound)
			return
		}
		defer rows.Close()

	} else if r.Method == http.MethodGet {

		// Fetch posts from the database
		rows, err = db.Query("SELECT id FROM Post ORDER BY created_at DESC")
		if err != nil {
			log.Println("Error fetching all posts:", err)
			ErrorHandler(w, "error2InHomePage", "error", http.StatusNotFound)
			return
		}
		defer rows.Close()

	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		post, err := getPostDetails(id)

		if err != nil {
			ErrorHandler(w, "Internal Server Error", "error", http.StatusInternalServerError)
		}
		data.Posts = append(data.Posts, *post)

	}

	RenderTemplate(w, "index", data)
}
