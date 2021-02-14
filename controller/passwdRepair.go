package controller

// import (
// 	"HypertubeAuth/controller/mailer"
// 	"HypertubeAuth/postgres"
// 	"HypertubeAuth/logger"
// 	"HypertubeAuth/model"
// 	"net/http"
// 	"strconv"
// )

// func passwdRepair(w http.ResponseWriter, r *http.Request) {
// 	email, Err := parseEmailFromRequest(r)
// 	if Err != nil {
// 		logger.Warning(r, Err.Error())
// 		errorResponse(w, Err)
// 		return
// 	}

// 	user, Err := postgres.UserGetBasicByEmail(email)
// 	if Err != nil {
// 		logger.Warning(r, Err.Error())
// 		errorResponse(w, Err)
// 		return
// 	}


// 	// Нужен не accessToken а код подтверждения почты
// 	// accessToken, Err := hash.CreateToken(user)
// 	// if Err != nil {
// 	// 	logger.Warning(r, "cannot get password hash - "+Err.Error())
// 	// 	errorResponse(w, Err)
// 	// 	return
// 	// }

// 	conf, Err := getConfig()
// 	if Err != nil {
// 		logger.Error(r, Err)
// 		errorResponse(w, Err)
// 		return
// 	}

// 	go func(user *model.UserBasic, serverPort uint) {
// 		if Err := mailer.SendEmailPasswdRepairMessage(user, serverPort); Err != nil {
// 			logger.Error(r, Err)
// 		} else {
// 			logger.Success(r, "Писмьмо для подтверждения почты пользователя #"+
// 				logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+" успешно отправлено")
// 		}
// 	}(user, conf.ServerPort)

// 	successResponse(w, nil)
// 	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
// 		" подал успешную заявку на восстановление пароля")
// }