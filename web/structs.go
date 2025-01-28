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

type CategoryDetails struct {
	CategoryID int
	Name       string
}

type PageDetails struct {
	LoggedIn         bool
	Posts            []PostDetails
	Categories       []CategoryDetails
	SelectedFilter   string
	SelectedCategory string
}
