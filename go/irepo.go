package main

type Post struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type Comment struct {
	Post
	Ref int `json:"ref"`
}

type PostColletion struct {
	Posts []Post `json:"posts"`
}

type CommentColletion struct {
	Comments []Comment `json:"comments"`
}

type IRepo interface {
	CreatePost(p Post) error
	CreateComment(cm Comment) error
	ReadAllPosts() (PostColletion, error)
	ReadPostComments(id int) (CommentColletion, error)
	Reset() error
}
