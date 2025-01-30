package web

import (
	"database/sql"
	"forum/database"
	"log"
	"strings"
)

// GetCategories retrieves all categories from the database
func GetCategories() ([]CategoryDetails, error) {
	rows, err := db.Query("SELECT id, name FROM Category")
	if err != nil {
		log.Println("Error retrieving categories:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryDetails
	for rows.Next() {
		var category CategoryDetails
		if err := rows.Scan(&category.CategoryID, &category.CategoryName); err != nil {
			log.Println("Error scanning category:", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
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
