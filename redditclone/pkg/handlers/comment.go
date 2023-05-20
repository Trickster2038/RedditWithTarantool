package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/comment"
	"redditclone/pkg/token"
	"time"

	"github.com/gorilla/mux"
)

func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["POST_ID"]

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		processReadRequestDataErr(w, err, h.Logger)
		return
	}

	var commentData map[string]string

	err = json.Unmarshal(reqBody, &commentData)
	if err != nil {
		processRequestBodyErr(w, reqBody, err, h.Logger)
		return
	}

	userID, username, err := token.GetClaims(r)
	if err != nil {
		processJwtClaimsErr(w, err, h.Logger)
		return
	}

	var cm comment.Comment
	cm.Body = commentData["comment"]
	cm.Created = time.Now().Format(time.RFC3339)
	cm.Author.ID = userID
	cm.Author.Username = username

	err = h.CommentRepo.Add(postID, cm)
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}

	postData, err := h.PostRepo.GetByID(postID)
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}
	err = h.getPostDependentData(&postData)
	if err != nil {
		processDependentDataErr(w, err, h.Logger)
		return
	}

	err = json.NewEncoder(w).Encode(postData)
	if err != nil {
		h.Logger.Error(ErrEncode)
	}
}

func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["POST_ID"]
	commentID := params["COMMENT_ID"]

	_, username, err := token.GetClaims(r)
	if err != nil {
		processJwtClaimsErr(w, err, h.Logger)
		return
	}

	err = h.CommentRepo.Delete(postID, commentID, username)
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}
	postData, err := h.PostRepo.GetByID(postID)
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}

	err = h.getPostDependentData(&postData)
	if err != nil {
		processDependentDataErr(w, err, h.Logger)
		return
	}

	err = json.NewEncoder(w).Encode(postData)
	if err != nil {
		h.Logger.Error(ErrEncode)
	}
}
