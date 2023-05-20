package user

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type UserMemoryRepository struct {
	data      map[string]User
	mu        sync.RWMutex
	idCounter int
}

func NewMockMemoryRepo() *UserMemoryRepository {
	repo := &UserMemoryRepository{
		data: make(map[string]User),
	}
	repo.Register("user", "12345678") //nolint: errcheck
	return repo
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: make(map[string]User),
	}
}

func (repo *UserMemoryRepository) Register(login, pass string) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, userExists := repo.data[login]
	if userExists {
		return user, ErrAlreadyExists
	} else {
		h := sha256.New()
		h.Write(append([]byte(pass), Salt...))
		repo.data[login] = User{repo.idCounter, login, fmt.Sprintf("%x", h.Sum(nil))}
		repo.idCounter++
	}
	return repo.data[login], nil
}

func (repo *UserMemoryRepository) Authorize(login, pass string) (User, error) {
	repo.mu.RLock()
	repo.mu.RUnlock()

	u, ok := repo.data[login]
	if !ok {
		return User{}, ErrNoUser
	}

	h := sha256.New()
	h.Write(append([]byte(pass), Salt...))
	if fmt.Sprintf("%x", h.Sum(nil)) != u.Password {
		return User{}, ErrBadPass
	}

	return u, nil
}
