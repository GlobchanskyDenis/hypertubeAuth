package controller

import (
	"HypertubeAuth/controller/hash"
	"HypertubeAuth/controller/validator"
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/postgres"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
**	/api/password/patch
**	Обновление пароля пользователя
**	В запросе должны содержаться поля passwd, new_passwd
**	авторизация в авторизационном хидере access_token
 */

func passwordPatch(w http.ResponseWriter, r *http.Request) {
	passwd, newPasswd, Err := parsePasswordsFromRequest(r)
	if Err != nil {
		logger.Warning(r, Err.Error())
		errorResponse(w, Err)
		return
	}

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

	if Err = validator.ValidatePassword(newPasswd); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	user, Err := postgres.UserGetBasicById(header.UserId)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	encryptedPass, Err := hash.PasswdHash(passwd)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if user.EncryptedPass == nil || *encryptedPass != *user.EncryptedPass {
		logger.Warning(r, "Хэши паролей не совпали. Ожидалось "+logger.BLUE+*user.EncryptedPass+logger.NO_COLOR+
			" получили "+logger.BLUE+*encryptedPass+logger.NO_COLOR)
		errorResponse(w, errors.ImpossibleToExecute.SetArgs("Пароль неверен", "Incorrect password"))
		return
	}

	user.EncryptedPass, Err = hash.PasswdHash(newPasswd)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if Err = postgres.UserUpdateEncryptedPassBasic(user); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	successResponse(w, nil)
	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" успешно обновил свои поля")
}

func parsePasswordsFromRequest(r *http.Request) (string, string, *errors.Error) {
	type Passwords struct {
		NewPasswd *string `json:"new_passwd"`
		Passwd    *string `json:"passwd"`
	}
	var pass = Passwords{}
	if err := json.NewDecoder(r.Body).Decode(&pass); err != nil {
		return "", "", errors.InvalidRequestBody.SetOrigin(err)
	}
	if pass.NewPasswd == nil {
		return "", "", errors.NoArgument.SetArgs("new_passwd", "new_passwd")
	}
	if pass.Passwd == nil {
		return "", "", errors.NoArgument.SetArgs("passwd", "passwd")
	}
	return *pass.Passwd, *pass.NewPasswd, nil
}
