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

/*
**	/api/email/resend
**	Повторная отправка кода подтверждения почты на почту
 */

func emailConfirm(w http.ResponseWriter, r *http.Request) {
	confirmCode, Err := parseCodeFromRequest(r)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	email, Err := hash.EmailHashDecode(confirmCode)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	user, Err := postgres.UserGetBasicByEmail(email)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if user.IsEmailConfirmed == true {
		logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
			" уже подтвердил свою почту")
		successResponse(w, nil)
		return
	}

	if user.EmailConfirmHash != confirmCode {
		logger.Warning(r, "Хэш подтверждения почты не совпал. Ожидалось "+user.EmailConfirmHash+" получено "+confirmCode)
		errorResponse(w, errors.ImpossibleToExecute.SetArgs("неверный код подтверждения", "confirm code is wrong"))
		return
	}

	if Err = postgres.UserConfirmEmailBasic(user); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	logger.Success(r, "Повторное письмо пользователя #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" обработано и поставлено в очередь на отправку")
	successResponse(w, nil)
}

func parseCodeFromRequest(r *http.Request) (string, *errors.Error) {
	var confirmCode = r.FormValue("code")
	if confirmCode != "" {
		return confirmCode, nil
	}

	if r.Body == nil {
		return "", errors.NoArgument.SetArgs("Отсутствует код подтверждения",
			"confirm code expected").SetHidden("Тело запроса пустое")
	}

	type Request struct {
		Code *string `json:"code"`
	}
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return "", errors.InvalidRequestBody.SetOrigin(err)
	}
	if request.Code == nil || *request.Code == "" {
		return "", errors.NoArgument.SetArgs("Отсутствует код подтверждения", "confirm code expected")
	}
	return *request.Code, nil
}
