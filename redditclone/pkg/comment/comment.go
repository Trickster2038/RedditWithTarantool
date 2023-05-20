package comment

import "errors"

type Comment struct {
	Created string `json:"created"`
	Author  struct {
		Username string `json:"username"`
		ID       string `json:"id"`
	} `json:"author"`
	Body string `json:"body"`
	ID   string `json:"id"`
}

var (
	ErrNoComment = errors.New("no comment found")
	ErrNoPost    = errors.New("no post found")
	ErrNoAccess  = errors.New("current user have no access for this action")
)

type CommentRepo interface {
	GetAllCommentsForPost(postID string) ([]Comment, error)
	Delete(postID string, commentID string, username string) error
	Add(postID string, cm Comment) error
}
