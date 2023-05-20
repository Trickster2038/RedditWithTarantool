package post

import (
	"encoding/json"
	"sort"
	"strconv"
	"sync"
)

type PostMemoryRepository struct {
	data      []Post
	mu        sync.RWMutex
	idCounter int
}

func NewMockMemoryRepo() *PostMemoryRepository {
	posts := make([]Post, 0)
	rawPosts := []string{
		`{"score":1,"views":0,"type":"link","title":"Link5Prog","author":{"username":"user","id":"0"},"category":"programming","url":"https://stackoverflow.com/questions/28999735/what-is-the-shortest-way-to-simply-sort-an-array-of-structs-by-arbitrary-field","votes":[{"user":"0","vote":1}],"Comments":[],"created":"2023-03-28T11:45:04+03:00","upvotePercentage":100,"id":"4"}`,
		`{"score":1,"views":0,"type":"link","title":"Link4Programming","author":{"username":"user","id":"0"},"category":"programming","url":"https://pkg.go.dev/sort#Slice","votes":[{"user":"0","vote":1}],"Comments":[],"created":"2023-03-28T11:44:48+03:00","upvotePercentage":100,"id":"3"}`,
		`{"score":1,"views":0,"type":"text","title":"Title3Prog","author":{"username":"user","id":"0"},"category":"programming","text":"programming is fun =)","votes":[{"user":"0","vote":1}],"Comments":[],"created":"2021-03-28T11:44:21+03:00","upvotePercentage":100,"id":"2"}`,
		`{"score":1,"views":0,"type":"text","title":"Title1Music","author":{"username":"user","id":"0"},"category":"music","text":"Music...","votes":[{"user":"0","vote":1}],"Comments":[],"created":"2023-03-28T11:42:14+03:00","upvotePercentage":100,"id":"0"}`,
		`{"score":1,"views":0,"type":"text","title":"Music2","author":{"username":"user","id":"0"},"category":"music","text":"Curt Cobain forever","votes":[{"user":"0","vote":1}],"Comments":[],"created":"2023-01-28T11:42:50+03:00","upvotePercentage":100,"id":"1"}`,
	}
	for _, v := range rawPosts {
		var post Post
		json.Unmarshal([]byte(v), &post) //nolint: errcheck
		posts = append(posts, post)
	}
	return &PostMemoryRepository{
		data:      posts,
		idCounter: len(posts),
	}
}

func NewMemoryRepo() *PostMemoryRepository {
	return &PostMemoryRepository{
		data: make([]Post, 0, 10),
	}
}

func (repo *PostMemoryRepository) GetAll() ([]Post, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return repo.data, nil
}

func (repo *PostMemoryRepository) GetAllInCategory(category string) ([]Post, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	result := make([]Post, 0)
	for _, v := range repo.data {
		if v.Category == category {
			result = append(result, v)
		}
	}
	return result, nil
}

func (repo *PostMemoryRepository) GetByID(id string) (Post, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, v := range repo.data {
		if v.ID == id {
			return v, nil
		}
	}
	return Post{}, ErrNoPost
}

func (repo *PostMemoryRepository) IncViewsByID(id string) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for i, _ := range repo.data {
		if repo.data[i].ID == id {
			repo.data[i].Views++
			return nil
		}
	}
	return ErrNoPost
}

func (repo *PostMemoryRepository) GetAllByUser(username string) ([]Post, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	result := make([]Post, 0)
	for _, v := range repo.data {
		if v.Author.Username == username {
			result = append(result, v)
		}
	}
	return result, nil
}

func (repo *PostMemoryRepository) Add(p Post) (int, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	p.ID = strconv.Itoa(repo.idCounter)
	repo.data = append(repo.data, p)
	repo.idCounter++
	return repo.idCounter - 1, nil
}

func (repo *PostMemoryRepository) Delete(postID string, username string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	found := false
	ps := make([]Post, 0)
	for _, v := range repo.data {
		if v.ID != postID {
			ps = append(ps, v)
		} else {
			if v.Author.Username != username {
				return ErrNoAccess
			}
			found = true
		}
	}

	if !found {
		return ErrNoPost
	} else {
		repo.data = ps
		return nil
	}
}

func SortByScoreDesc(sl []Post) []Post {
	sort.SliceStable(sl, func(i, j int) bool {
		return sl[i].Score > sl[j].Score // ">" - for desc order
	})
	return sl
}
