package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ErrBody         = errors.New("unable to read request body")
	ErrUnmarshal    = errors.New("unable to unmarhal request JSON")
	ErrDBConnection = errors.New("unable to connect to DB")
	ErrQueryExec    = errors.New("query execution error")
)

const (
	MsgOk = "ok"
)

func wrappedCreatePostHandler(repo IRepo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("createPostHandler")

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SafeJsonEncode(w, map[string]string{"error": ErrBody.Error()})
		}

		var p Post
		err = json.Unmarshal(reqBody, &p)
		if err != nil {
			SafeJsonEncode(w, map[string]string{"error": ErrUnmarshal.Error()})
		}

		err = repo.CreatePost(p)

		if err == nil {
			SafeJsonEncode(w, map[string]string{"status": MsgOk})
		} else {
			SafeJsonEncode(w, map[string]string{"error": ErrQueryExec.Error()})
		}
	}
}
