package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrBody        = errors.New("unable to read request body")
	ErrUnmarshal   = errors.New("unable to unmarhal request JSON")
	ErrQueryParams = errors.New("not all url params found")
)

const (
	MsgOk = "ok"
)

func wrappedCreatePostHandler(repo IRepo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("createPostHandler")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			SafeJsonEncode(w, map[string]string{"error": ErrBody.Error()})
			return
		}

		var p Post
		err = json.Unmarshal(reqBody, &p)
		if err != nil {
			SafeJsonEncode(w, map[string]string{"error": ErrUnmarshal.Error()})
			return
		}

		err = repo.CreatePost(p)

		if err == nil {
			SafeJsonEncode(w, map[string]string{"status": MsgOk})
		} else {
			SafeJsonEncode(w, map[string]string{"error": err.Error()})
		}
	}
}

func wrappedCreateCommentHandler(repo IRepo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("createCommentHandler")

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			SafeJsonEncode(w, map[string]string{"error": ErrQueryParams.Error()})
			return
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			SafeJsonEncode(w, map[string]string{"error": ErrBody.Error()})
			return
		}

		var cm Comment
		err = json.Unmarshal(reqBody, &cm)
		if err != nil {
			SafeJsonEncode(w, map[string]string{"error": ErrUnmarshal.Error()})
			return
		}

		err = repo.CreateComment(id, cm)
		if err == nil {
			SafeJsonEncode(w, map[string]string{"status": MsgOk})
		} else {
			SafeJsonEncode(w, map[string]string{"error": err.Error()})
		}
	}
}

func wrappedReadAllPostsHandler(repo IRepo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("readAllPostsHandler")

		res, err := repo.ReadAllPosts()
		if err == nil {
			SafeJsonEncode(w, res)
		} else {
			SafeJsonEncode(w, map[string]string{"error": ErrQueryExec.Error()})
		}
	}
}

func wrappedReadPostCommentsHandler(repo IRepo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("readPostComments")

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			SafeJsonEncode(w, map[string]string{"error": ErrQueryExec.Error()})
			return
		}

		res, err := repo.ReadPostComments(id)
		if err == nil {
			SafeJsonEncode(w, res)
		} else {
			SafeJsonEncode(w, map[string]string{"error": ErrQueryExec.Error()})
		}
	}
}

func wrappedResetHandler(repo IRepo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("resetHandler")

		err := repo.Reset()

		if err == nil {
			SafeJsonEncode(w, map[string]string{"status": MsgOk})
		} else {
			SafeJsonEncode(w, map[string]string{"error": err.Error()})
		}
	}
}
