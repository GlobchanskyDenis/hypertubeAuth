package controller

import (
	"HypertubeAuth/controller/hash"
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/model"
	"encoding/json"
	"net/http"
	"strconv"
)

func userProfile(w http.ResponseWriter, r *http.Request) {
	var token model.Token
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		logger.Warning(r, errors.InvalidRequestBody.SetOrigin(err).Error())
		errorResponse(w, errors.InvalidRequestBody)
		return
	}

	user, Err := hash.GetUserFromToken(token.AccessToken)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		logger.Error(r, errors.MarshalError.SetOrigin(err))
		errorResponse(w, errors.MarshalError)
		return
	}

	successResponse(w, jsonUser)
	logger.Success(r, "user #" + strconv.Itoa(int(user.UserId)) + " was checked")
}