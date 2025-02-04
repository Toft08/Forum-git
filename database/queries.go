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
			COALESCE(likes.post_likes, 0) AS post_likes,
   			COALESCE(likes.post_dislikes, 0) AS post_dislikes,
			COALESCE(GROUP_CONCAT(Category.name, ','), '') AS categories
		FROM Post
		LEFT JOIN User ON Post.user_id = User.id
		LEFT JOIN (
			SELECT 
				post_id,
				SUM(CASE WHEN type = 1 THEN 1 ELSE 0 END) AS post_likes,
				SUM(CASE WHEN type = 2 THEN 1 ELSE 0 END) AS post_dislikes
			FROM Like
			GROUP BY post_id
		) AS likes ON Post.id = likes.post_id
		LEFT JOIN Post_Category ON Post.id = Post_Category.post_id
		LEFT JOIN Category ON Post_Category.category_id = Category.id
		WHERE Post.id = ?
		GROUP BY Post.id, Post.user_id, User.username, Post.title, Post.content, Post.created_at;
	`
	return query
}

// CommentContent returns the query to fetch comment details
func CommentContent() string {
	query := `
		SELECT 
			Comment.id AS comment_id,
			Comment.post_id AS post_id,
			Comment.content AS comment_content,
			Comment.user_id,
			Comment.created_at,
			User.username AS username,
			COUNT(CASE WHEN Like.type = 1 THEN 1 END) AS comment_likes,
			COUNT(CASE WHEN Like.type = 2 THEN 1 END) AS comment_dislikes
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

// Votes returns the query to fetch Like.type for a post or a comment
func Votes() string {
	query := `
    SELECT type
    FROM Like
    WHERE user_id = ?
      AND (post_id = COALESCE(?, post_id) AND comment_id = COALESCE(?, comment_id));
	`
	return query
}
