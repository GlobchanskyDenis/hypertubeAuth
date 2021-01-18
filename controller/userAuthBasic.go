package controller

import (
	"HypertubeAuth/controller/hash"
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/postgres"
	"net/http"
	"strconv"
)

func userAuthBasic(w http.ResponseWriter, r *http.Request) {
	email, passwd, ok := r.BasicAuth()
	if !ok {
		logger.Warning(r, "authenticaion failed - email or password not found")
		errorResponse(w, errors.NoArgument.SetArgs("Отсутствует авторизационное поле", "Authorization field expected"))
		return
	}

	user, Err := postgres.UserGetBasicByEmail(email)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	encryptedPass, Err := hash.PasswdHash(passwd)
	if Err != nil {
		logger.Warning(r, "cannot get password hash - " + Err.Error())
		errorResponse(w, Err)
		return
	}

	if user.EncryptedPass != encryptedPass {
		logger.Warning(r, "authenticaion failed - password missmatch")
		errorResponse(w, errors.AuthFail)
		return
	}

	if user.IsEmailConfirmed == false {
		logger.Warning(r, "authenticaion failed - email of user is not confirmed")
		errorResponse(w, errors.NotConfirmedMail)
		return
	}

	tokenRaw, Err := hash.CreateToken(user.TransformToUser(), "user_basic")
	if Err != nil {
		logger.Warning(r, "cannot get password hash - " + Err.Error())
		errorResponse(w, Err)
		return
	}
	successResponse(w, tokenRaw)
	logger.Success(r, "user #" + strconv.Itoa(int(user.UserId)) + " was authenticated")
}
