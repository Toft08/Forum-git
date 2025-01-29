package web

import (
	"database/sql"
	"forum/database"
	"log"
	"net/http"
	"strconv"
)

func HomePage(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	switch r.Method {
	case http.MethodGet:
		HandleHomeGet(w, r, data)
	case http.MethodPost:
		HandleHomePost(w, r, data)
	default:
		ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func HandleHomeGet(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	data.LoggedIn, _ = VerifySession(r)
	// Fetch posts from the database
	rows, err := db.Query("SELECT id FROM Post ORDER BY created_at DESC")
	if err != nil {
		// log.Println("Error fetching posts:", err)
		ErrorHandler(w, "error2InHomePage", http.StatusNotFound)
		return
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

func HandleHomePost(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	var args []interface{}
	var userID int
	var rows *sql.Rows
	var err error
	var query string

	data.LoggedIn, userID = VerifySession(r)

	data.SelectedCategory = r.FormValue("topic")
	categoryID, err := strconv.Atoi(data.SelectedCategory)
	if err != nil {
		log.Println("Error converting categoryID", err)
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
	}
	if categoryID > 0 {
		query = database.FilterCategories()
	}

	data.LoggedIn, _ = VerifySession(r)
	if data.LoggedIn {
		args = append(args, userID)
		data.SelectedFilter = r.FormValue("filter")

		switch data.SelectedFilter {
		case "createdByMe":
			query = "SELECT p.id FROM Post p WHERE p.user_id = ?"
		case "likedByMe":
			query = database.MyLikes()
		case "dislikedByMe":
			query = database.MyDislikes()
		}

		if categoryID != 0 {
			query += " JOIN Post_category pc ON p.id = pc.post_id AND pc.category_id = ?"
			args = append(args, categoryID)

		}
		query += " ORDER BY p.created_at DESC"
		// Fetch posts from the database for a specific user
		rows, err = db.Query(query, args...)
		if err != nil {
			log.Println("Error fetching posts by filter:", err)
			ErrorHandler(w, "errorFetchingPosts", http.StatusNotFound)
			return
		}
	}

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
