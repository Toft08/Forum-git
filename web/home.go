package web

import (
	"database/sql"
	"forum/database"
	"log"
	"net/http"
	"strconv"
)

// HomePage handles the rendering of the home page
func HomePage(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	data.ValidationError = ""
	switch r.Method {
	case http.MethodGet:
		HandleHomeGet(w, r, data)
	case http.MethodPost:
		HandleHomePost(w, r, data)
	default:
		ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// HandleHomeGet fetches posts from the database and renders the home page
func HandleHomeGet(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	data.LoggedIn, _ = VerifySession(r)
	// Fetch posts from the database
	rows, err := db.Query(`
        SELECT Post.id
        FROM Post
        ORDER BY Post.created_at DESC;
    `)
	// rows, err := db.Query(`
	//     SELECT p.id, p.title, p.content, u.username
	//     FROM Post p
	//     JOIN User u ON p.user_id = u.id
	//     ORDER BY p.created_at DESC;
	// `)
	if err != nil {
		log.Println("Error fetching posts:", err)
		ErrorHandler(w, "error2InHomePage", http.StatusNotFound)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		//var title, content, username string
		if err := rows.Scan(&id); err != nil {
			ErrorHandler(w, "Error scanning post ID", http.StatusInternalServerError)
			return
		}
		post, err := GetPostDetails(id, 0)

		if err != nil {
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		}
		data.Posts = append(data.Posts, *post)

	}

	RenderTemplate(w, "index", data)
}

// HandleHomePost handles the filtering of posts based on the user's selection
func HandleHomePost(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	var args []interface{}
	var userID int
	var rows *sql.Rows
	var err error
	var query string
	var categoryID int

	data.LoggedIn, userID = VerifySession(r)
	data.SelectedFilter = r.FormValue("filter")
	data.SelectedCategory = r.FormValue("topic")

	if !data.LoggedIn && data.SelectedFilter != "" {
		log.Println("User not logged in")
		return
	}

	if data.LoggedIn {
		if data.SelectedCategory == "" && data.SelectedFilter == "" {
			HandleHomeGet(w, r, data)
			return
		} else if data.SelectedCategory != "" && data.SelectedFilter == "" {
			categoryID, err = strconv.Atoi(data.SelectedCategory)
			if err != nil {
				log.Println("Error converting categoryID", err)
				ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			}
			query = database.FilterCategories()
			args = append(args, categoryID)
		} else {
			args = append(args, userID)
			switch data.SelectedFilter {
			case "createdByMe":
				query = "SELECT Post.id FROM Post WHERE Post.user_id = ?"
			case "likedByMe":
				query = database.MyLikes()
			case "dislikedByMe":
				query = database.MyDislikes()
			}

		}
	} else {
		if data.SelectedCategory == "" {
			HandleHomeGet(w, r, data)
			return
		} else {
			categoryID, err = strconv.Atoi(data.SelectedCategory)
			if err != nil {
				log.Println("Error converting categoryID", err)
				ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			}
			query = database.FilterCategories()
			args = append(args, categoryID)
		}
	}
	query += " ORDER BY Post.created_at DESC;"
	// Fetch posts from the database for a specific user
	rows, err = db.Query(query, args...)
	if err != nil {
		log.Println("Error fetching posts by filter:", err)
		ErrorHandler(w, "errorFetchingPosts", http.StatusNotFound)
		return
	}

	for rows.Next() {
		var id int
		rows.Scan(&id)
		post, err := GetPostDetails(id, userID)

		if err != nil {
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		}
		data.Posts = append(data.Posts, *post)

	}

	RenderTemplate(w, "index", data)
}
