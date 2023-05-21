package main

import (
	"context"

	tarantool "github.com/viciious/go-tarantool"
)

type TarantoolRepo struct {
	Host     string
	User     string
	Password string
}

func (r *TarantoolRepo) connect() (*tarantool.Connection, error) {
	opts := tarantool.Options{User: r.User, Password: r.Password}
	conn, err := tarantool.Connect(r.Host, &opts)
	if err != nil {
		return &tarantool.Connection{}, err
	}
	return conn, nil
}

func (r *TarantoolRepo) CreatePost(p Post) error {
	conn, err := connect(host, user, pass)
	defer conn.Close()
	if err != nil {
		return err
	}

	query := &tarantool.Eval{
		Expression: "box.space.post:auto_increment{...}",
		Tuple:      []interface{}{p.Content},
	}
	resp := conn.Exec(context.Background(), query)
	defer conn.Close()

	return resp.Error
}
