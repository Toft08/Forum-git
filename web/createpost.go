package web

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

// CreatePost receives details for created post and inserts them into the database
func CreatePost(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	var userID int
	var err error
	var categoryIDs []int

	data.LoggedIn, _, data.Username = VerifySession(r)

	if r.Method == http.MethodPost {
		data.LoggedIn, userID, data.Username = VerifySession(r)
		if !data.LoggedIn {
			ErrorHandler(w, "You must be logged in to create a post", http.StatusUnauthorized)
			return
		}

		err = r.ParseForm()
		if err != nil {
			log.Println("Unable to parse form:", err)
			ErrorHandler(w, "Bad Request", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["category"]

		if title == "" || content == "" {
			ErrorHandler(w, "Title or content cannot be empty", http.StatusBadRequest)
			return
		}

		if len(categories) == 0 {
			categories = append(categories, "1") // If no category chosen, give category id 1 (=general)
		}

		// Converting categoryIDs to integers and validating that they exists in the database
		for _, cat := range categories {
			var categoryID int
			categoryID, err = HandleCategory(cat)
			if err != nil {
				log.Println("Error handling categoryID in createpost", err)
				ErrorHandler(w, "Bad Request", http.StatusBadRequest)
			}

			categoryIDs = append(categoryIDs, categoryID)
		}

		err = AddPostToDatabase(title, content, categoryIDs, userID)
		if err != nil {
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)

	} else if r.Method != http.MethodGet {

		ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	RenderTemplate(w, "create-post", data)
}

// AddPostToDatabase inserts a new post into the database
func AddPostToDatabase(title, content string, categories []int, userID int) error {

	var result sql.Result
	var err error
	result, err = db.Exec("INSERT INTO Post (title, content, user_id, created_at) VALUES (?, ?, ?, ?)",
		title, content, userID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error inserting post:", err)
		return err
	}

	// Get the post id for the post inserted
	postID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting post ID:", err)
		return err
	}

	// Add all categories into Post_category table
	for _, categoryID := range categories {
		_, err = db.Exec("INSERT INTO Post_category (category_id, post_id) VALUES (?, ?)",
			categoryID, postID)
		if err != nil {
			log.Println("Error inserting post category:", err)
			return err
		}
	}

	return nil
}
