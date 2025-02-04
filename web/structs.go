package web

type CommentDetails struct {
	CommentID   int
	PostID      int
	Content     string
	UserID      int
	Username    string
	CreatedAt   string
	Likes       int
	Dislikes    int
	LikedNow    bool
	DislikedNow bool
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
	LikedNow    bool
	DislikedNow bool
}

type PageDetails struct {
	LoggedIn         bool
	Username         string
	Categories       []CategoryDetails
	Posts            []PostDetails
	SelectedCategory string
	SelectedFilter   string
	ValidationError  string
}

type CategoryDetails struct {
	CategoryID   int
	CategoryName string
}
