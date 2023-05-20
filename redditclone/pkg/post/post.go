package post

import (
	"errors"
	"redditclone/pkg/comment"
	"redditclone/pkg/vote"
)

type Post struct {
	Score  int    `json:"score"`
	Views  int    `json:"views"`
	Type   string `json:"type"`
	Title  string `json:"title"`
	Author struct {
		Username string `json:"username"`
		ID       string `json:"id"`
	} `json:"author"`
	Category         string            `json:"category"`
	Text             string            `json:"text,omitempty"` // optional
	URL              string            `json:"url,omitempty"`  // optional
	Votes            []vote.Vote       `json:"votes"`
	Comments         []comment.Comment `json:"comments"`
	Created          string            `json:"created"`
	UpvotePercentage int               `json:"upvotePercentage"`
	ID               string            `json:"id"`
}

var (
	ErrNoAccess = errors.New("current user have no access for this action")
	ErrNoPost   = errors.New("post does not exist")
)

func (p *Post) UpdateStats() {
	p.Score = 0
	upvotes := 0
	for _, v := range p.Votes {
		p.Score += v.Vote
		if v.Vote == 1 {
			upvotes++
		}
	}
	if len(p.Votes) > 0 {
		p.UpvotePercentage = upvotes * 100 / len(p.Votes)
	} else {
		p.UpvotePercentage = 0
	}
}

type PostRepo interface {
	GetAll() ([]Post, error)
	GetAllInCategory(category string) ([]Post, error)
	GetAllByUser(username string) ([]Post, error)
	GetByID(ID string) (Post, error)
	IncViewsByID(ID string) error
	Add(ps Post) (int, error)
	Delete(postID string, username string) error
}
