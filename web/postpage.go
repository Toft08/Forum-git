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
	var err error
	data.LoggedIn, userID = VerifySession(r)

	if !data.LoggedIn {
		ErrorHandler(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
		return
	}
	vote := r.FormValue("vote")
	commentID := r.FormValue("comment-id")
	content := r.FormValue("comment")

	if content != "" {
		// Insert comment into the database
		_, err = db.Exec("INSERT INTO Comment (post_id, content, user_id, created_at) VALUES (?, ?, ?, ?)",
			postID, content, userID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error creating post:", err)
			ErrorHandler(w, "errorInCreatePost", http.StatusNotFound)
			return
		}
	} else {
		var likeType int
		var post int
		var comment int
		if vote == "like" {
			likeType = 1
		} else if vote == "dislike" {
			likeType = 2
		} else {
			log.Println("Invalid vote value: ", vote)
			ErrorHandler(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if commentID == "" {
			comment = 0
			post = postID
		} else {
			comment, err = strconv.Atoi(commentID)
			if err != nil {
				log.Println("Error converting commentID", err)
				ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			exists := ValidateCommentID(comment)
			if !exists {
				log.Println("CommentID doesn't exist", comment)
				ErrorHandler(w, "Bad Request", http.StatusBadRequest)
				return
			}
			post = 0
		}
		err = AddVotes(userID, post, comment, likeType)
		if err != nil {
			log.Printf("Error adding votes to the database: userID %d, postID %d, commentID %d, like type %d\n", userID, post, comment, likeType)
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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

// ValidateCommentID checks if a comment with the given ID exists in the database
func ValidateCommentID(commentID int) bool {
	var comment int
	err := db.QueryRow("SELECT id FROM Comment WHERE id = ?", commentID).Scan(&comment)
	if err != nil {
		log.Println("Error scanning commentID:", err)
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
		updateQuery := `UPDATE Like SET type = 0, created_at = ? WHERE user_id = ? AND (post_id = ? OR comment_id = ?)`
		_, err = db.Exec(updateQuery, time.Now().Format("2006-01-02 15:04:05"), userID, postID, commentID)
		if err != nil {
			log.Println("Error updating the record to 0:", err)
			return err
		}
	} else if likeType == -1 {
		// If no record exists, insert a new one
		insertQuery := `INSERT INTO Like (type, user_id, post_id, comment_id, created_at) VALUES (?, ?, ?, ?, ?)`
		_, err = db.Exec(insertQuery, vote, userID, postID, commentID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error inserting record:", err)
			return err
		}
	} else {
		updateQuery := `UPDATE Like SET type = ?, created_at = ? WHERE user_id = ? AND (post_id = ? OR comment_id = ?)`
		_, err = db.Exec(updateQuery, vote, time.Now().Format("2006-01-02 15:04:05"), userID, postID, commentID)
		if err != nil {
			log.Println("Error updating the record to new vote:", err)
			return err
		}
	}
	return nil
}
