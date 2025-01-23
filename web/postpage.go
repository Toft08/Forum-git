package web

import (
	"forum/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CommentDetails struct {
	CommentID int
	Content   string
	UserID    int
	Username  string
	Likes     int
	Dislikes  int
}

type PostDetails struct {
	PostID      int
	UserID      int
	Username    string
	PostTitle   string
	PostContent string
	Comments    []CommentDetails
	Categories  []string
	CreatedAt   string
	Likes       int
	Dislikes    int
}

// PostHandler handles requests to view a specific post
func PostHandler(w http.ResponseWriter, r *http.Request) {

	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		log.Println("Error converting postID to int:", err)
		errorHandler(w, "Page Not Found", "error", http.StatusNotFound)
	}

	post, err := getPostDetails(postID)
	if err != nil {
		log.Println("Error fetching post details:", err)
		errorHandler(w, "Internal Server Error", "error", http.StatusInternalServerError)
		return
	}

	comments, err := getComments(postID)
	if err != nil {
		log.Println("Error fetching comments:", err)
		errorHandler(w, "Internal Server Error", "error", http.StatusInternalServerError)
		return
	}

	log.Println(comments)
	post.Comments = comments

	renderTemplate(w, "post", post)

}

// getPostDetails fetches the details of a specific post from the database
func getPostDetails(postID int) (*PostDetails, error) {

	row := db.QueryRow(database.PostContent(), postID)

	// Scan the data into a PostDetails struct
	post := PostDetails{}
	//var categories string
	err := row.Scan(
		&post.PostID,
		&post.UserID,
		&post.Username,
		&post.PostTitle,
		&post.PostContent,
		&post.CreatedAt,
		//&categories,
		&post.Likes,
		&post.Dislikes,
	)
	if err != nil {
		return nil, err
	}

	// Split the concatenated categories into a slice
	// DATABASE MISSING CAT POST CONNECTION
	// post.Categories = []string{}
	// if categories != "" {
	// 	catSlice := strings.Fields(categories)
	// 	post.Categories = catSlice

	// }
	return &post, nil
}

// getComments fetches all comments for a specific post from the database
func getComments(postID int) ([]CommentDetails, error) {

	rows, err := db.Query(database.CommentContent(), postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CommentDetails
	for rows.Next() {
		var comment CommentDetails
		err := rows.Scan(
			&comment.CommentID,
			&comment.Content,
			&comment.UserID,
			&comment.Username,
			&comment.Likes,
			&comment.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
