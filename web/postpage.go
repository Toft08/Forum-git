package web

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PostHandler handles requests to view a specific post
func PostHandler(w http.ResponseWriter, r *http.Request, data *PageDetails) {

	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		log.Println("Error converting postID to int:", err)
		ErrorHandler(w, "Page Not Found", http.StatusNotFound)
	}

	valid := ValidatePostID(postID)
	if !valid {
		log.Println("Invalid postID")
		ErrorHandler(w, "Page Not Found", http.StatusNotFound)
	}

	switch r.Method {
	case http.MethodGet:
		HandlePostPageGet(w, r, data, postID)
	case http.MethodPost:
		HandlePostPagePost(w, r, data, postID)
	default:
		ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

// HandlePostPageGet handles get requests to the post page
func HandlePostPageGet(w http.ResponseWriter, r *http.Request, data *PageDetails, postID int) {
	var userID int
	data.LoggedIn, userID = VerifySession(r)
	data.Posts = nil

	post, err := GetPostDetails(postID, userID)
	if err != nil {
		log.Println("Error fetching post details:", err)
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Posts = append(data.Posts, *post)

	RenderTemplate(w, "post", data)
}

// HandlePostPagePost handles post requests to the post page
func HandlePostPagePost(w http.ResponseWriter, r *http.Request, data *PageDetails, postID int) {
	var userID int
	data.LoggedIn, userID = VerifySession(r)

	if !data.LoggedIn {
		ErrorHandler(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
		return
	}

	content := r.FormValue("comment")

	// Insert comment into the database
	_, err := db.Exec("INSERT INTO Comment (post_id, content, user_id, created_at) VALUES (?, ?, ?, ?)",
		postID, content, userID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error creating post:", err)
		ErrorHandler(w, "errorInCreatePost", http.StatusNotFound)
		return
	}
	HandlePostPageGet(w, r, data, postID)
}

// ValidatePostID checks if a post with the given ID exists in the database
func ValidatePostID(postID int) bool {
	var post int
	err := db.QueryRow("SELECT id FROM Post WHERE id = ?", postID).Scan(&post)
	if err != nil {
		log.Println("Error scanning postID:", err)
		return false
	}
	return true
}

// AddVotes adds or updates a vote type for a post or comment
func AddVotes(userID, postID, commentID, vote int) error {

	if postID == 0 && commentID == 0 {
		return fmt.Errorf("both postID and commentID cannot be zero")
	}

	query := `SELECT Type FROM Like WHERE user_id = ? AND (post_id = ? OR comment_id = ?)`
	row := db.QueryRow(query, userID, postID, commentID)
	var likeType int
	err := row.Scan(&likeType)
	if err != nil {
		if err == sql.ErrNoRows {
			likeType = -1 // To imply that no record exists
		} else {
			log.Println("Error scanning current value:", err)
			return err
		}
	}
	if likeType == vote {
		// If existing like type is the same the the current, remove the like by changing the type to 0
		updateQuery := `UPDATE Like SET type = 0 WHERE user_id = ? AND (post_id = ? OR comment_id = ?)`
		_, err = db.Exec(updateQuery, userID, postID, commentID)
		if err != nil {
			log.Println("Error updating the record to 0:", err)
			return err
		}
	} else if likeType == -1 {
		// If no record exists, insert a new one
		insertQuery := `INSERT INTO Like (type, user_id, post_id, comment_id) VALUES (?, ?, ?, ?)`
		_, err = db.Exec(insertQuery, vote, userID, postID, commentID)
		if err != nil {
			log.Println("Error inserting record:", err)
			return err
		}
	} else {
		updateQuery := `UPDATE Like SET type = ? WHERE user_id = ? AND (post_id = ? OR comment_id = ?)`
		_, err = db.Exec(updateQuery, vote, userID, postID, commentID)
		if err != nil {
			log.Println("Error updating the record to new vote:", err)
			return err
		}
	}
	return nil
}
