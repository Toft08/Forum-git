package web

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"
)

// CreatePost receives details for created post and inserts them into the database
func CreatePost(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	var userID int
	var err error

	data.LoggedIn, userID = VerifySession(r)

	if r.Method == http.MethodPost {
		if !data.LoggedIn {
			ErrorHandler(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
			// http.Error(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		categories := r.Form["category"]

		log.Println("Received title:", title)
		log.Println("Received content:", content)

		// Insert post into the database
		var result sql.Result
		result, err = db.Exec("INSERT INTO Post (title, content, user_id, created_at) VALUES (?, ?, ?, ?)",
			title, content, userID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			ErrorHandler(w, "errorInCreatePost", http.StatusNotFound)
			return
		}

		// Get the post id for the post inserted
		postID, err := result.LastInsertId()
		if err != nil {
			ErrorHandler(w, "errorGettingPostID", http.StatusInternalServerError)
			return
		}
		// If no category chosen, give category id 1 (=general)
		if len(categories) == 0 {
			categories = append(categories, "1")
		}
		// Add all categories into Post_category table
		for _, cat := range categories {
			var categoryID int
			categoryID, err = strconv.Atoi(cat)
			if err != nil {
				log.Println("Error converting categoryID", err)
				ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)

			}
			_, err = db.Exec("INSERT INTO Post_category (category_id, post_id) VALUES (?, ?)",
				categoryID, postID)
			if err != nil {
				ErrorHandler(w, "errorInInsertCategory", http.StatusNotFound)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)

	} else if r.Method != http.MethodGet {

		ErrorHandler(w, "Wrong method", http.StatusMethodNotAllowed)
	}

	RenderTemplate(w, "create-post", data)
}
