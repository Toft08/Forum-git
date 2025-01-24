package web

import (
	"forum/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// PostHandler handles requests to view a specific post
func PostHandler(w http.ResponseWriter, r *http.Request) {

	data := PageDetails{}

	IsLoggedIn, _ := IsLoggedIn(r)

	data.LoggedIn = IsLoggedIn

	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		log.Println("Error converting postID to int:", err)
		ErrorHandler(w, "Page Not Found", "error", http.StatusNotFound)
	}

	post, err := getPostDetails(postID)
	if err != nil {
		log.Println("Error fetching post details:", err)
		ErrorHandler(w, "Internal Server Error", "error", http.StatusInternalServerError)
		return
	}

	data.Posts = append(data.Posts, *post)

	RenderTemplate(w, "post", data)

}

// getPostDetails fetches the details of a specific post from the database
func getPostDetails(postID int) (*PostDetails, error) {

	row := db.QueryRow(database.PostContent(), postID)
	log.Println("Query row")
	var err error
	// Scan the data into a PostDetails struct
	post := PostDetails{}
	var categories string
	log.Println("starting scanning")
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
		log.Println(err)
		return nil, err
	}
	log.Println("scanning done")

	if categories != "" {
		post.Categories = strings.Split(categories, ",")
	}

	log.Println("Adding comments")
	postComments, err := getComments(postID)
	if err != nil {
		log.Println("Error getting comments")
		return nil, err
	}
	post.Comments = postComments

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
