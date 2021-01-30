package mailer

import (
	"HypertubeAuth/errors"
	"HypertubeAuth/model"
	"net/smtp"
	"strconv"
)

type Config struct {
	Host   string `conf:"host"`
	Email  string `conf:"email"`
	Passwd string `conf:"passwd"`
}

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	return cfg
}

func getConfig() (*Config, *errors.Error) {
	if cfg == nil {
		return nil, errors.NotConfiguredPackage.SetArgs("controller/mailer", "controller/mailer")
	}
	return cfg, nil
}

func SendEmailConfirmMessage(user *model.UserBasic, serverPort uint) *errors.Error {
	conf, Err := getConfig()
	if Err != nil {
		return Err
	}

	portString := strconv.FormatUint(uint64(serverPort), 10)

	auth := smtp.PlainAuth("", conf.Email, conf.Passwd, conf.Host)
	message := `To: <` + user.Email + `>
From: "Hypertube administration" <` + conf.Email + `>
Subject: Confirm email in project Hypertube
MIME-Version: 1.0
Content-type: text/html; charset=utf8

<html><head></head><body>
<span style="font-size: 1.3em; color: green;">Hello, ` + user.Username + `, click below to confirm your email
<form method="POST" action="http://localhost:`+portString+`/user/update/status/">
	<input type="hidden" name="x-reg-token" value="`+user.EmailConfirmHash+`">
	<input type="submit" value="Click to confirm mail">
</form>
<a target="_blank" href="http://localhost:`+portString+`/user/update/status/?x-reg-token=`+user.EmailConfirmHash+`">click to confirm mail</a></br> 
 `+user.EmailConfirmHash+`</br>
if this letter came by mistake - delete it 
</span></body></html>
`

	if err := smtp.SendMail(conf.Host+":587", auth, conf.Email, []string{user.Email}, []byte(message)); err != nil {
		return errors.MailerError.SetOrigin(err)
	}
	return nil
}