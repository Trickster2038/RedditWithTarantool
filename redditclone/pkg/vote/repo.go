package vote

import (
	"sync"
)

type VoteMemoryRepository struct {
	data map[string][]Vote
	mu   sync.RWMutex
}

func NewMockMemoryRepo() *VoteMemoryRepository {
	vts := make(map[string][]Vote)
	vts["0"] = []Vote{{"0", 1}}
	vts["1"] = []Vote{{"0", 1}}
	vts["2"] = []Vote{{"0", 1}}
	vts["3"] = []Vote{{"0", 1}}
	vts["4"] = []Vote{{"0", 1}}
	return &VoteMemoryRepository{
		data: vts,
	}
}

func NewMemoryRepo() *VoteMemoryRepository {
	return &VoteMemoryRepository{
		data: make(map[string][]Vote),
	}
}

func (repo *VoteMemoryRepository) GetAllVotesForPost(
	postID string,
) (
	[]Vote,
	error,
) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	res, ok := repo.data[postID]
	if (!ok) || (res == nil) {
		return []Vote{}, nil
	} else {
		return res, nil
	}
}

func (repo *VoteMemoryRepository) Vote(postID string, userID string, voteVal int) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	vts, ok := repo.data[postID]
	if (!ok) || (vts == nil) {
		vts = []Vote{}
	}

	// delete old vote (unvote)
	vtBuf := make([]Vote, 0)
	for _, v := range vts {
		if v.User != userID {
			vtBuf = append(vtBuf, v)
		}
	}
	repo.data[postID] = vtBuf

	if voteVal != 0 { // vote
		repo.data[postID] = append(repo.data[postID], Vote{userID, voteVal})
	}

	return nil
}
