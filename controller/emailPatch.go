package controller

import (
	"HypertubeAuth/controller/hash"
	"HypertubeAuth/controller/mailer"
	"HypertubeAuth/logger"
	"HypertubeAuth/model"
	"HypertubeAuth/postgres"
	"net/http"
	"strconv"
)

/*
**	/api/email/patch
**	изменение почты. Отправляет письмо для подтверждения на новую почту
**	Первый эндпоинт из двух. Первый дергает пользователь с сайта,
**	второй - с почты для подтверждения
**	-- еще не оттестировано !!!!!!!
 */
func emailPatch(w http.ResponseWriter, r *http.Request) {
	email, Err := parseEmailFromRequest(r)
	if Err != nil {
		logger.Warning(r, Err.Error())
		errorResponse(w, Err)
		return
	}

	user, Err := postgres.UserGetBasicByEmail(email)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	user.NewEmail = &email

	user.EmailConfirmHash, Err = hash.EmailHashEncode(email)
	if Err != nil {
		logger.Warning(r, "cannot get password hash - "+Err.Error())
		errorResponse(w, Err)
		return
	}

	if Err = postgres.UserUpdateBasic(user); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	go func(user *model.UserBasic, serverIp string, serverPort uint) {
		if Err := mailer.SendEmailPatchMailAddress(user, serverIp, serverPort); Err != nil {
			logger.Error(r, Err)
		} else {
			logger.Success(r, "Писмьмо для подтверждения новой почты пользователя #"+
				logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+" успешно отправлено")
		}
	}(user, conf.ServerIp, conf.ServerPort)

	successResponse(w, nil)
	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" подал успешную заявку на изменение почтового адреса")
}
