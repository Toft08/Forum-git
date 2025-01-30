package web

import (
	"database/sql"
	"forum/database"
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

	// Check the cookie and get userID
	var userID int
	data.LoggedIn, userID = VerifySession(r)
	data.Posts = nil

	if r.Method == http.MethodPost {
		data.LoggedIn, _ = VerifySession(r)

		if !data.LoggedIn {
			ErrorHandler(w, "Unauthorized: You must be logged in to create a post", http.StatusUnauthorized)
			return
		}

		content := r.FormValue("comment")

		// Insert post into the database
		_, err := db.Exec("INSERT INTO Comment (post_id, content, user_id, created_at) VALUES (?, ?, ?, ?)",
			postID, content, userID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error creating post:", err)
			ErrorHandler(w, "errorInCreatePost", http.StatusNotFound)
			return
		}
	} else if r.Method != http.MethodGet {
		ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	post, err := GetPostDetails(postID)
	if err != nil {
		log.Println("Error fetching post details:", err)
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Posts = append(data.Posts, *post)

	RenderTemplate(w, "post", data)

}

// GetPostDetails fetches the details of a specific post from the database
func GetPostDetails(postID int) (*PostDetails, error) {

	row := db.QueryRow(database.PostContent(), postID)
	var err error
	// Scan the data into a PostDetails struct
	post := PostDetails{}
	var categories string
	err = row.Scan(
		&post.PostID,
		&post.UserID,
		&post.Username,
		&post.PostTitle,
		&post.PostContent,
		&post.CreatedAt,
		&post.Likes,
		&post.Dislikes,
		&categories,
	)

	if err != nil {
		log.Println("Error scanning rows")
		return nil, err
	}

	if categories != "" {
		post.Categories = strings.Split(categories, ",")
	}

	post.Comments, err = GetComments(postID)
	if err != nil {
		log.Println("Error getting comments")
		return nil, err
	}

	post.Votes, err = GetVotes(postID, 0)
	if err != nil {
		log.Println("Error getting votes")
		return nil, err
	}
	return &post, nil
}

// GetComments fetches all comments for a specific post from the database
func GetComments(postID int) ([]CommentDetails, error) {

	rows, err := db.Query(database.CommentContent(), postID)
	if err != nil {
		log.Println("Error fetching comments from database")
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
			log.Println("Error scanning rows")
			return nil, err
		}
		comment.Votes, err = GetVotes(0, comment.CommentID)
		if err != nil {
			log.Println("Error getting votes")
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// GetVotes fetches the votes for a specific post or comment from the database
func GetVotes(postID, commentID int) (map[int]int, error) {
	var rows *sql.Rows
	var err error
	if postID == 0 {
		rows, err = db.Query(database.PostVotes(), nil, commentID)
	} else {
		rows, err = db.Query(database.PostVotes(), postID, nil)
	}
	if err != nil {
		log.Println("Error fetching votes from database")
		return nil, err
	}
	defer rows.Close()

	votes := make(map[int]int)
	for rows.Next() {
		var userID, vote int
		err := rows.Scan(&userID, &vote)
		if err != nil {
			log.Println("Error scanning rows")
			return nil, err
		}
		votes[vote] = userID
	}

	return votes, nil
}
