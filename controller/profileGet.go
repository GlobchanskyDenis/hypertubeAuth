package controller

import (
	"HypertubeAuth/controller/hash"
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/postgres"
	"encoding/json"
	"net/http"
	"strconv"
)

func profileGet(w http.ResponseWriter, r *http.Request) {
	/*
	**	Получаю токен из заголовка
	*/
	accessToken := r.Header.Get("access_token")
	if accessToken == "" {
		logger.Error(r, errors.UserNotLogged.SetArgs("отсутствует токен доступа", "access token expected"))
		errorResponse(w, errors.UserNotLogged)
		return
	}

	header, Err := hash.GetHeaderFromToken(accessToken)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	var id uint
	var idString = r.FormValue("user_id")
	if idString != "" {
		idInt, err := strconv.Atoi(idString)
		if err != nil {
			logger.Warning(r, "Id пользователя ("+idString+") содежит ошибку.")
			errorResponse(w, errors.InvalidRequestBody)
			return
		}
		id = uint(idInt)
	} else {
		id = header.UserId
	}

	user, Err := postgres.UserGetBasicById(id)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	/*
	**	Если пользователь не мой - почистить приватные поля
	*/
	if user.UserId != header.UserId {
		user.Sanitize()
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