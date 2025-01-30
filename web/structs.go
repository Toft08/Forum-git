package web

type CommentDetails struct {
	CommentID int
	Content   string
	UserID    int
	Username  string
	Likes     int
	Dislikes  int
	Votes     map[int]int
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
	Votes       map[int]int
}

type PageDetails struct {
	LoggedIn         bool
	Categories       []CategoryDetails
	Posts            []PostDetails
	SelectedCategory string
	SelectedFilter   string
}

type CategoryDetails struct {
	CategoryID   int
	CategoryName string
}
