package controller

import (
	"HypertubeAuth/controller/hash"
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/postgres"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

/*
**	/api/email/confirm
**	Подтверждение почтового адреса
**	-- Проверено
 */
func emailConfirm(w http.ResponseWriter, r *http.Request) {
	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())), //base64.StdEncoding.EncodeToString(Err.ToJson())
			http.StatusTemporaryRedirect)
		return
	}

	confirmCode, Err := parseCodeFromRequest(r)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	email, Err := hash.EmailHashDecode(confirmCode)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	user, Err := postgres.UserGetBasicByEmail(email)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	if user.IsEmailConfirmed == true {
		logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
			" уже подтвердил свою почту")
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect,
			http.StatusTemporaryRedirect)
		return
	}

	// if user.EmailConfirmHash != confirmCode {
	// 	logger.Warning(r, "Хэш подтверждения почты не совпал. Ожидалось "+user.EmailConfirmHash+" получено "+confirmCode)
	// 	Err := errors.ImpossibleToExecute.SetArgs("неверный код подтверждения", "confirm code is wrong")
	// 	http.Redirect(w, r,
	// 		conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
	// 		http.StatusTemporaryRedirect)
	// 	return
	// }

	if Err = postgres.UserConfirmEmailBasic(user); Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" подтвердил свой почтовый адрес "+logger.BLUE+user.Email+logger.NO_COLOR)
	http.Redirect(w, r,
		conf.SocketRedirect+conf.ErrorRedirect,
		http.StatusTemporaryRedirect)
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
