package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"redditclone/pkg/weberrors"

	"go.uber.org/zap"
)

var (
	ErrEncode       = errors.New("JSON encoding error")
	errEncodeErrMsg = errors.New("error description encoding error")
)

func processRepoErr(w http.ResponseWriter, err error, logger *zap.SugaredLogger) {
	logger.Warn(err.Error())
	errStruct := weberrors.DetailedErrors{
		Errors: []weberrors.DetailedError{{
			Location: "repo",
			Param:    "",
			Value:    "",
			Msg:      err.Error(),
		},
		}}
	w.WriteHeader(http.StatusInternalServerError)
	err = json.NewEncoder(w).Encode(errStruct)
	if err != nil {
		logger.Error(errEncodeErrMsg)
	}
}

func processRequestBodyErr(w http.ResponseWriter, reqBody []byte, err error, logger *zap.SugaredLogger) {
	logger.Warn(err.Error())
	errStruct := weberrors.DetailedErrors{
		Errors: []weberrors.DetailedError{{
			Location: "body",
			Param:    "request body",
			Value:    string(reqBody),
			Msg:      err.Error(),
		},
		}}
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(errStruct)
	if err != nil {
		logger.Error(errEncodeErrMsg)
	}
}

func processJwtClaimsErr(w http.ResponseWriter, err error, logger *zap.SugaredLogger) {
	logger.Warn(err.Error())
	errStruct := weberrors.DetailedErrors{
		Errors: []weberrors.DetailedError{{
			Location: "JWT-token claims",
			Param:    "",
			Value:    "",
			Msg:      err.Error(),
		},
		}}
	w.WriteHeader(http.StatusInternalServerError)
	err = json.NewEncoder(w).Encode(errStruct)
	if err != nil {
		logger.Error(errEncodeErrMsg)
	}
}

func processDependentDataErr(w http.ResponseWriter, err error, logger *zap.SugaredLogger) {
	logger.Warn(err.Error())
	errStruct := weberrors.DetailedErrors{
		Errors: []weberrors.DetailedError{{
			Location: "Dependent post data",
			Param:    "",
			Value:    "",
			Msg:      err.Error(),
		},
		}}
	w.WriteHeader(http.StatusInternalServerError)
	err = json.NewEncoder(w).Encode(errStruct)
	if err != nil {
		logger.Error(errEncodeErrMsg)
	}
}

func processReadRequestDataErr(w http.ResponseWriter, err error, logger *zap.SugaredLogger) {
	logger.Warn(err.Error())
	errStruct := weberrors.DetailedErrors{
		Errors: []weberrors.DetailedError{{
			Location: "request data",
			Param:    "",
			Value:    "",
			Msg:      err.Error(),
		},
		}}
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(errStruct)
	if err != nil {
		logger.Error(errEncodeErrMsg)
	}
}
