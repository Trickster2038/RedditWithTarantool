package user

import "errors"

type User struct {
	ID       int
	Username string `json:"username"`
	Password string // unmarshal plain password, store hash
}

var (
	ErrNoUser        = errors.New("user not found")
	ErrBadPass       = errors.New("invald password")
	ErrAlreadyExists = errors.New("already exists")
	Salt             = []byte("password salt =)")
)

type UserRepo interface {
	Register(login, pass string) (User, error)
	Authorize(login, pass string) (User, error)
}
