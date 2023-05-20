package vote

type Vote struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

type VoteRepo interface {
	GetAllVotesForPost(postID string) ([]Vote, error)
	Vote(postID string, userID string, voteVal int) error
}
