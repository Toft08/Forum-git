package web

import (
	"log"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	if r.URL.Path != "/" {
		ErrorHandler(w, "error1InHomepage", http.StatusNotFound)
		return
	}
	var err error
	data.LoggedIn, _, err = VerifySession(r)
	if err != nil {
		log.Println("Error verifying session:", err)
	}

	//message := r.URL.Query().Get("message")

	// Fetch posts from the database
	rows, err := db.Query("SELECT id FROM Post ORDER BY created_at DESC")
	if err != nil {
		// log.Println("Error fetching posts:", err)
		ErrorHandler(w, "error2InHomePage", http.StatusNotFound)
		return
	}
	defer rows.Close()

	for rows.Next() {
		log.Println("Post id rows")
		var id int
		rows.Scan(&id)
		post, err := getPostDetails(id)

		if err != nil {
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		}
		data.Posts = append(data.Posts, *post)

	}
	// Pass posts to template
	RenderTemplate(w, "index", data)
}
