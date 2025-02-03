package database

// PostContent returns the query to fetch post details
func PostContent() string {
	query := `
		SELECT 
			Post.id AS post_id,
			Post.user_id AS user_id,
			User.username AS username, 
			Post.title AS post_title,
			Post.content AS post_content,
			Post.created_at AS post_created_at,
			COUNT(CASE WHEN Like.type = 0 THEN 1 END) AS post_likes,
			COUNT(CASE WHEN Like.type = 1 THEN 1 END) AS post_dislikes,
			COALESCE(GROUP_CONCAT(Category.name, ','), '') AS categories
		FROM post
		LEFT JOIN user ON Post.user_id = User.id
		LEFT JOIN like ON Post.id = Like.post_id
		LEFT JOIN Post_Category ON Post.id = Post_Category.post_id
		LEFT JOIN Category ON Post_Category.category_id = Category.id
		WHERE Post.id = ?
		GROUP BY Post.id;

	`
	return query
}

// CommentContent returns the query to fetch comment details
func CommentContent() string {
	query := `
		SELECT 
			Comment.id AS comment_id,
			Comment.content AS comment_content,
			Comment.user_id,
			User.username AS username,
			COUNT(CASE WHEN Like.type = 0 THEN 1 END) AS comment_likes,
			COUNT(CASE WHEN Like.type = 1 THEN 1 END) AS comment_dislikes
		FROM comment
		LEFT JOIN user ON Comment.user_id = User.id
		LEFT JOIN like ON Comment.id = Like.comment_id
		WHERE Comment.post_id = ?
		GROUP BY Comment.id, User.id;
`
	return query
}

// MyLikes returns the query to fetch posts liked by the user
func MyLikes() string {
	query := `
	SELECT
		Post.id 
	FROM Post 
	INNER JOIN Like ON Post.id = Like.post_id
	WHERE Like.user_id = ? AND Like.type = 1
	`

	return query
}

// MyDislikes returns the query to fetch posts disliked by the user
func MyDislikes() string {
	query := `
	SELECT
		Post.id 
	FROM Post 
	INNER JOIN Like ON Post.id = Like.post_id
	WHERE Like.user_id = ? AND Like.type = 2
	`

	return query
}
// FilterCategories returns the query to filter posts by category
func FilterCategories() string {
	query := `    
	SELECT Post.id
	FROM Post
	JOIN Post_category ON Post.id = Post_category.post_id
	WHERE Post_category.category_id = ?	
	`
	return query

}
// PostComments returns the query to fetch comments for a post
func PostVotes() string {
	query := `
		SELECT user_id, type
		FROM Like
		WHERE type IN (1, 2) AND post_id = ?;
	`
	return query
}
// CommentVotes returns the query to fetch votes for a comment
func Votes() string {
	query := `
    SELECT userID, type
    FROM "Like"
    WHERE type IN (1, 2)
      AND (postID = COALESCE(?, postID) AND commentID = COALESCE(?, commentID));
`
	return query
}
