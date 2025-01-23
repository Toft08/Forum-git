package database

// PostContent returns the query to fetch post details
func PostContent() string {
	query := `
		SELECT 
			post.id AS post_id,
			post.user_id AS user_id,
			user.username AS username, 
			post.title AS post_title,
			post.content AS post_content,
			post.created_at AS post_created_at,
			COUNT(CASE WHEN like.type = 0 THEN 1 END) AS post_likes,
			COUNT(CASE WHEN like.type = 1 THEN 1 END) AS post_dislikes
		FROM post
		LEFT JOIN user ON post.user_id = user.id
		LEFT JOIN like ON post.id = like.post_id
		WHERE post.id = ?
		GROUP BY post.id;
	`
	return query
}

// CommentContent returns the query to fetch comment details
func CommentContent() string {
	query := `
		SELECT 
			comment.id AS comment_id,
			comment.content AS comment_content,
			comment.user_id,
			user.username AS username,
			COUNT(CASE WHEN like.type = 0 THEN 1 END) AS comment_likes,
			COUNT(CASE WHEN like.type = 1 THEN 1 END) AS comment_dislikes
		FROM comment
		LEFT JOIN user ON comment.user_id = user.id
		LEFT JOIN like ON comment.id = like.comment_id
		WHERE comment.post_id = ?
		GROUP BY comment.id, user.id;
`
	return query
}
