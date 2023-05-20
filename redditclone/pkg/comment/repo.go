package comment

import (
	"strconv"
	"sync"
)

type CommentMemoryRepository struct {
	data      map[string][]Comment
	mu        sync.Mutex
	idCounter uint32
}

func NewMockMemoryRepo() *CommentMemoryRepository {
	return &CommentMemoryRepository{
		data: make(map[string][]Comment),
	}
}

func NewMemoryRepo() *CommentMemoryRepository {
	return &CommentMemoryRepository{
		data: make(map[string][]Comment),
	}
}

func (repo *CommentMemoryRepository) Add(postID string, cm Comment) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	cm.ID = strconv.Itoa(int(repo.idCounter))
	repo.idCounter++
	repo.data[postID] = append(repo.data[postID], cm)
	return nil
}

func (repo *CommentMemoryRepository) Delete(postID string, commentID string, username string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	cms := make([]Comment, 0)
	currentCms, ok := repo.data[postID]
	if ok {
		found := false
		for _, v := range currentCms {
			if v.ID != commentID {
				cms = append(cms, v)
			} else {
				if v.Author.Username != username {
					return ErrNoAccess
				}
				found = true
			}
		}
		if found {
			repo.data[postID] = cms
			return nil
		} else {
			return ErrNoComment
		}
	} else {
		return ErrNoPost
	}
}

func (repo *CommentMemoryRepository) GetAllCommentsForPost(
	postID string,
) (
	[]Comment,
	error,
) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	res, ok := repo.data[postID]
	if (!ok) || (res == nil) {
		return []Comment{}, nil
	} else {
		return res, nil
	}
}
