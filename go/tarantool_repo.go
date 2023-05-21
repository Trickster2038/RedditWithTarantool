package main

import (
	"context"
	"errors"

	tarantool "github.com/viciious/go-tarantool"
)

type TarantoolRepo struct {
	Host     string
	User     string
	Password string
}

var (
	ErrNoPost     = errors.New("no such post found")
	ErrConnection = errors.New("unable to connect to DB")
	ErrQueryExec  = errors.New("query execution error")
)

func (r *TarantoolRepo) connect() (*tarantool.Connection, error) {
	opts := tarantool.Options{User: r.User, Password: r.Password}
	conn, err := tarantool.Connect(r.Host, &opts)
	if err != nil {
		return &tarantool.Connection{}, ErrConnection
	}
	return conn, nil
}

func (repo *TarantoolRepo) CreatePost(p Post) error {
	conn, err := repo.connect()

	if err != nil {
		return ErrConnection
	}

	query := &tarantool.Eval{
		Expression: "box.space.post:auto_increment{...}",
		Tuple:      []interface{}{p.Content},
	}
	resp := conn.Exec(context.Background(), query)
	defer conn.Close()

	if resp.Error != nil {
		return ErrQueryExec
	}

	return nil
}

func (repo *TarantoolRepo) CreateComment(cm Comment) error {
	conn, err := repo.connect()
	if err != nil {
		return ErrConnection
	}
	defer conn.Close()

	query := &tarantool.Select{Space: "post", Index: "primary", Key: cm.Ref}
	resp := conn.Exec(context.Background(), query)

	if resp.Error != nil || len(resp.Data) == 0 {
		return ErrNoPost
	}

	query2 := &tarantool.Eval{
		Expression: "box.space.comment:auto_increment{...}",
		Tuple:      []interface{}{cm.Content, cm.Ref},
	}
	resp = conn.Exec(context.Background(), query2)

	if resp.Error != nil {
		return ErrQueryExec
	}

	return nil
}

func (repo *TarantoolRepo) ReadAllPosts() (PostColletion, error) {
	conn, err := repo.connect()
	if err != nil {
		return PostColletion{}, ErrConnection
	}
	defer conn.Close()

	query := &tarantool.Select{Space: "post"}
	resp := conn.Exec(context.Background(), query)

	if resp.Error != nil {
		return PostColletion{}, ErrQueryExec
	}

	result := PostColletion{make([]Post, 0)}
	p := Post{}
	for _, tuple := range resp.Data {
		p.ID = int(tuple[0].(int64))
		p.Content = tuple[1].(string)
		result.Posts = append(result.Posts, p)
	}

	return result, nil
}

func (repo *TarantoolRepo) ReadPostComments(id int) (CommentColletion, error) {
	conn, err := repo.connect()
	if err != nil {
		return CommentColletion{}, ErrConnection
	}
	defer conn.Close()

	query := &tarantool.Select{Space: "comment", Index: "ref_idx", Key: id}
	resp := conn.Exec(context.Background(), query)

	if resp.Error != nil {
		return CommentColletion{}, ErrQueryExec
	}

	result := CommentColletion{make([]Comment, 0)}
	cm := Comment{}
	for _, tuple := range resp.Data {
		cm.ID = int(tuple[0].(int64))
		cm.Content = tuple[1].(string)
		cm.Ref = int(tuple[2].(int64))
		result.Comments = append(result.Comments, cm)
	}

	return result, nil
}

func (repo *TarantoolRepo) Reset() error {
	conn, err := repo.connect()
	if err != nil {
		return ErrConnection
	}
	defer conn.Close()

	query := &tarantool.Eval{Expression: "box.space.post:truncate()"}
	resp := conn.Exec(context.Background(), query)
	if resp.Error != nil {
		return ErrQueryExec
	}

	query = &tarantool.Eval{Expression: "box.space.comment:truncate()"}
	resp = conn.Exec(context.Background(), query)
	if resp.Error != nil {
		return ErrQueryExec
	}

	return nil
}
