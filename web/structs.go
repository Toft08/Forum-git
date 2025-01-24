package web

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
