package handlers

import (
	"encoding/json"
	"net/http"
	"redditclone/pkg/token"

	"github.com/gorilla/mux"
)

func (h *PostsHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, 1)
}

func (h *PostsHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, -1)
}

func (h *PostsHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, 0)
}

func (h *PostsHandler) vote(w http.ResponseWriter, r *http.Request, voteVal int) {
	params := mux.Vars(r)
	postID := params["POST_ID"]

	userID, _, err := token.GetClaims(r)
	if err != nil {
		processJwtClaimsErr(w, err, h.Logger)
		return
	}

	err = h.VoteRepo.Vote(postID, userID, voteVal)
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
