package web

import (
	"database/sql"
	"forum/database"
	"log"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	if r.URL.Path != "/" {
		ErrorHandler(w, "error1InHomepage", http.StatusNotFound)
		return
	}

	var userID int
	var rows *sql.Rows
	var err error
	query := "SELECT id FROM Post ORDER BY created_at DESC"

	data.LoggedIn, userID = VerifySession(r)

	if r.Method == http.MethodPost {
		data.SelectedCategory = r.FormValue("topic")
		if data.LoggedIn {
			data.SelectedFilter = r.FormValue("filter")

			switch data.SelectedFilter {
			case "createdByMe":
				query = "SELECT id FROM Post WHERE user_id = ? ORDER BY created_at DESC"
			case "likedByMe":
				query = database.MyLikes()
			case "dislikedByMe":
				query = database.MyDislikes()
			}
			// Fetch posts from the database for a specific user
			rows, err = db.Query(query, userID)
			if err != nil {
				log.Println("Error fetching posts by filter:", err)
				ErrorHandler(w, "errorFetchingPosts", http.StatusNotFound)
				return
			}
		}
	} else if r.Method == http.MethodGet {

		// Fetch posts from the database
		rows, err = db.Query(query)
		if err != nil {
			// log.Println("Error fetching posts:", err)
			ErrorHandler(w, "error2InHomePage", http.StatusNotFound)
			return
		}
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		rows.Scan(&id)
		post, err := getPostDetails(id)

		if err != nil {
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		}
		data.Posts = append(data.Posts, *post)

	}

	RenderTemplate(w, "index", data)
}
