package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	tarantool "github.com/viciious/go-tarantool"
)

type Post struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type Comment struct {
	Post
	Ref int `json:"ref"`
}

type PostsColletion struct {
	Posts []Post `json:"posts"`
}

type CommentColletion struct {
	Comments []Comment `json:"comments"`
}

var (
	host       = "127.0.0.1:3301"
	user       = "admin"
	pass       = "pass"
	accessPort = "8085"
)

func connect(host, user, pass string) (*tarantool.Connection, error) {
	opts := tarantool.Options{User: user, Password: pass}
	conn, err := tarantool.Connect(host, &opts)
	if err != nil {
		return &tarantool.Connection{}, err
	}
	return conn, nil
}

func readAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("readAllPostsHandler")
	conn, _ := connect(host, user, pass)
	defer conn.Close()

	query := &tarantool.Select{Space: "post"}
	resp := conn.Exec(context.Background(), query)

	if resp.Error != nil {
		w.Write([]byte(fmt.Sprintf("Select failed: %s", resp.Error)))
		return
	}

	payload := PostsColletion{make([]Post, 0)}
	p := Post{}
	for _, tuple := range resp.Data {
		p.ID = int(tuple[0].(int64))
		p.Content = tuple[1].(string)
		payload.Posts = append(payload.Posts, p)
	}

	json.NewEncoder(w).Encode(payload)
}

func readPostComments(w http.ResponseWriter, r *http.Request) {
	log.Println("readPostComments")

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	conn, _ := connect(host, user, pass)
	defer conn.Close()

	query := &tarantool.Select{Space: "comment", Index: "ref_idx", Key: id}
	resp := conn.Exec(context.Background(), query)

	if resp.Error != nil {
		w.Write([]byte(fmt.Sprintf("Select failed: %s", resp.Error)))
		return
	}

	payload := CommentColletion{make([]Comment, 0)}
	cm := Comment{}
	for _, tuple := range resp.Data {
		cm.ID = int(tuple[0].(int64))
		cm.Content = tuple[1].(string)
		cm.Ref = int(tuple[2].(int64))
		payload.Comments = append(payload.Comments, cm)
	}

	json.NewEncoder(w).Encode(payload)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("createPostHandler")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var p Post
	json.Unmarshal(reqBody, &p)

	conn, _ := connect(host, user, pass)
	defer conn.Close()

	query := &tarantool.Eval{
		Expression: "box.space.post:auto_increment{...}",
		Tuple:      []interface{}{p.Content},
	}
	resp := conn.Exec(context.Background(), query)

	if resp.Error == nil {
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte(fmt.Sprintf("%v", resp)))
	}
}

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("createCommentHandler")

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	conn, _ := connect(host, user, pass)
	defer conn.Close()

	query := &tarantool.Select{Space: "post", Index: "primary", Key: id}
	resp := conn.Exec(context.Background(), query)

	if resp.Error != nil {
		w.Write([]byte(fmt.Sprintf("Select failed: %s", resp.Error)))
		return
	}

	if len(resp.Data) == 0 {
		w.Write([]byte("No such post"))
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var cm Comment
	json.Unmarshal(reqBody, &cm)

	query2 := &tarantool.Eval{
		Expression: "box.space.comment:auto_increment{...}",
		Tuple:      []interface{}{cm.Content, cm.Ref},
	}
	resp = conn.Exec(context.Background(), query2)

	if resp.Error == nil {
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte(fmt.Sprintf("%v", resp)))
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/post", createPostHandler).Methods("Post")
	r.HandleFunc("/comment", createCommentHandler).Methods("Post")
	r.HandleFunc("/posts", readAllPostsHandler).Methods("Get")
	r.HandleFunc("/comments", readPostComments).Methods("Get")
	http.Handle("/", r)

	fmt.Printf("Server is listening on %s\n", accessPort)
	http.ListenAndServe(":"+accessPort, nil)
}
