package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	tarantool "github.com/viciious/go-tarantool"
)

type Post struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

var (
	host = "127.0.0.1:3301"
	user = "admin"
	pass = "pass"
)

func connect(host, user, pass string) (*tarantool.Connection, error) {
	opts := tarantool.Options{User: user, Password: pass}
	conn, err := tarantool.Connect(host, &opts)
	if err != nil {
		return &tarantool.Connection{}, err
	}
	return conn, nil
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("createPostHandler")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var p Post

	json.Unmarshal(reqBody, &p)

	conn, _ := connect(host, user, pass)
	query := &tarantool.Eval{
		Expression: "box.space.post:auto_increment{...}",
		Tuple:      []interface{}{p.Content},
	}

	resp := conn.Exec(context.Background(), query)
	log.Println(resp)
	if resp.Error == nil {
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte(fmt.Sprintf("%v", resp)))
	}
}

func main() {
	fmt.Println("(~_~) Hello (~_~)")
	r := mux.NewRouter()
	// r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/post", createPostHandler).Methods("Post")
	http.Handle("/", r)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8085", nil)
}
