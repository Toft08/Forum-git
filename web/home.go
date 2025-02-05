package web

import (
	"database/sql"
	"fmt"
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
		ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// HandleHomeGet fetches posts from the database and renders the home page
func HandleHomeGet(w http.ResponseWriter, r *http.Request, data *PageDetails) {
	data.LoggedIn, _, data.Username = VerifySession(r)

	// Fetch posts from the database
	rows, err := db.Query(`
        SELECT Post.id
        FROM Post
        ORDER BY Post.created_at DESC;
    `)
	if err != nil {
		log.Println("Error fetching posts:", err)
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
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

	data.LoggedIn, userID, data.Username = VerifySession(r)
	data.SelectedFilter = r.FormValue("filter")
	selectedCat := r.FormValue("topic")
	data.SelectedCategory = selectedCat

	if !data.LoggedIn && data.SelectedFilter != "" {
		log.Println("User not logged in")
		return
	}

	if data.LoggedIn {
		if data.SelectedCategory == "" && data.SelectedFilter == "" {
			HandleHomeGet(w, r, data)
			return
		} else if data.SelectedCategory != "" && data.SelectedFilter == "" {
			categoryID, err = HandleCategory(selectedCat)
			if err != nil {
				log.Println("Error handling category", err)
				ErrorHandler(w, "Bad request", http.StatusBadRequest)
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
			categoryID, err = HandleCategory(selectedCat)
			if err != nil {
				log.Println("Error handling category", err)
				ErrorHandler(w, "Bad request", http.StatusBadRequest)
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
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
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

// HandleCategory converts the category ID into a string and returns validated ID
func HandleCategory(category string) (int, error) {

	categoryID, err := strconv.Atoi(category)
	if err != nil {
		log.Println("Error converting categoryID", err)
		return 0, err
	}

	valid := ValidateCategoryID(categoryID)
	if !valid {
		log.Println("Invalid categoryID", category)
		return 0, fmt.Errorf("invalid category id: %s", category)
	}

	return categoryID, nil

}

// ValidateCategoryID checks if the category ID given exists in the databse
func ValidateCategoryID(categoryID int) bool {
	var category int
	err := db.QueryRow("SELECT id FROM Category WHERE id = ?", categoryID).Scan(&category)
	if err != nil {
		log.Println("Error scanning category ID:", err)
		return false
	}
	return true
}
