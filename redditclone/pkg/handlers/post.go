package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/comment"
	"redditclone/pkg/post"
	"redditclone/pkg/token"
	"redditclone/pkg/vote"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostsHandler struct {
	Logger      *zap.SugaredLogger
	PostRepo    post.PostRepo
	CommentRepo comment.CommentRepo
	VoteRepo    vote.VoteRepo
}

func (h *PostsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	postsArr, err := h.PostRepo.GetAll()
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}

	err = h.getPostsDependentData(postsArr)
	if err != nil {
		processDependentDataErr(w, err, h.Logger)
		return
	}

	err = json.NewEncoder(w).Encode(postsArr)
	if err != nil {
		h.Logger.Error(ErrEncode)
	}
}

func (h *PostsHandler) GetAllInCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postsArr, err := h.PostRepo.GetAllInCategory(params["CATEGORY_NAME"])
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}

	err = h.getPostsDependentData(postsArr)
	if err != nil {
		processDependentDataErr(w, err, h.Logger)
	}

	err = json.NewEncoder(w).Encode(postsArr)
	if err != nil {
		h.Logger.Error(ErrEncode)
	}
}

func (h *PostsHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postsArr, err := h.PostRepo.GetAllByUser(params["USER_LOGIN"])
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}

	err = h.getPostsDependentData(postsArr)
	if err != nil {
		processDependentDataErr(w, err, h.Logger)
	}
	err = json.NewEncoder(w).Encode(postsArr)
	if err != nil {
		h.Logger.Error(ErrEncode)
	}
}

func (h *PostsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postData, err := h.PostRepo.GetByID(params["POST_ID"])
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}
	err = h.PostRepo.IncViewsByID(params["POST_ID"])
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

func (h *PostsHandler) Add(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		processReadRequestDataErr(w, err, h.Logger)
		return
	}

	var postData post.Post

	err = json.Unmarshal(reqBody, &postData)
	if err != nil {
		processRequestBodyErr(w, reqBody, err, h.Logger)
		return
	}

	userID, username, err := token.GetClaims(r)
	if err != nil {
		processRequestBodyErr(w, reqBody, err, h.Logger)
		return
	}

	postData.Score = 1
	postData.Author.ID = fmt.Sprintf("%v", userID)
	postData.Author.Username = fmt.Sprintf("%v", username)
	postData.Votes = make([]vote.Vote, 0)
	postData.Votes = append(postData.Votes, vote.Vote{User: userID, Vote: 1})
	postData.Comments = make([]comment.Comment, 0)
	postData.UpvotePercentage = 100
	postData.Created = time.Now().Format(time.RFC3339)

	postID, err := h.PostRepo.Add(postData)
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}

	err = h.VoteRepo.Vote(strconv.Itoa(postID), userID, 1)
	if err != nil {
		processRepoErr(w, err, h.Logger)
	}
	err = h.getPostDependentData(&postData)
	if err != nil {
		processDependentDataErr(w, err, h.Logger)
	}

	err = json.NewEncoder(w).Encode(postData)
	if err != nil {
		h.Logger.Error(ErrEncode)
	}
}

func (h *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	_, username, err := token.GetClaims(r)
	if err != nil {
		processJwtClaimsErr(w, err, h.Logger)
		return
	}

	err = h.PostRepo.Delete(params["POST_ID"], username)
	if err != nil {
		processRepoErr(w, err, h.Logger)
		return
	}

	payload := map[string]string{
		"message": "success",
	}

	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		h.Logger.Error(ErrEncode)
	}
}
