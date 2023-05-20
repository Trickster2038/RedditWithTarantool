package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/token"
	"redditclone/pkg/user"
	"redditclone/pkg/weberrors"
	"strconv"

	"go.uber.org/zap"
)

type UserHandler struct {
	Logger   *zap.SugaredLogger
	UserRepo user.UserRepo
}

type Token struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		processReadRequestDataErr(w, err, h.Logger)
		return
	}

	var userData user.User

	err = json.Unmarshal(reqBody, &userData)
	if err != nil {
		processRequestBodyErr(w, reqBody, err, h.Logger)
		return
	}

	userData, err = h.UserRepo.Register(userData.Username, userData.Password)
	if err != nil {
		h.Logger.Warnf(err.Error())
		errStruct := weberrors.DetailedErrors{
			Errors: []weberrors.DetailedError{{
				Location: "auth",
				Param:    "username",
				Value:    userData.Username,
				Msg:      err.Error(),
			},
			}}
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(errStruct)
		if err != nil {
			h.Logger.Error(ErrEncode)
		}
		return
	}

	token, err := token.GetToken(strconv.Itoa(userData.ID), userData.Username)
	if err != nil {
		h.Logger.Warnf(err.Error())
		errStruct := weberrors.DetailedErrors{
			Errors: []weberrors.DetailedError{{
				Location: "JWT-auth gen",
				Param:    "user",
				Value:    userData,
				Msg:      err.Error(),
			},
			}}
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(errStruct)
		if err != nil {
			h.Logger.Error(ErrEncode)
		}
		return
	}

	resp := map[string]interface{}{
		"token": token,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.Logger.Error(ErrEncode)
	}
}

func (h *UserHandler) Authorize(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		processReadRequestDataErr(w, err, h.Logger)
		return
	}

	var userData user.User

	err = json.Unmarshal(reqBody, &userData)
	if err != nil {
		processRequestBodyErr(w, reqBody, err, h.Logger)
		return
	}

	userData, err = h.UserRepo.Authorize(userData.Username, userData.Password)
	if err != nil {
		resp := map[string]interface{}{
			"message": err.Error(),
		}
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.Logger.Error(ErrEncode)
		}
		return
	}

	token, err := token.GetToken(strconv.Itoa(userData.ID), userData.Username)
	if err != nil {
		h.Logger.Warn(err.Error())
		errStruct := weberrors.DetailedErrors{
			Errors: []weberrors.DetailedError{{
				Location: "JWT-token gen",
				Param:    "username",
				Value:    userData.Username,
				Msg:      err.Error(),
			},
			}}
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(errStruct)
		if err != nil {
			h.Logger.Error(ErrEncode)
		}
		return
	}

	resp := map[string]interface{}{
		"token": token,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.Logger.Error(ErrEncode)
	}
}
