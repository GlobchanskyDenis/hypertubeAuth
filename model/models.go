package model

import (
	"HypertubeAuth/errors"
	"time"
)

type UserBasicModel struct {
	UserId           uint    `json:"user_id"`
	User42Id		 *uint	 `json:"-"`
	ImageBody        *string `json:"image_body"`
	Email            string  `json:"email"`
	EncryptedPass    *string `json:"-"`
	Fname            *string `json:"first_name"`
	Lname            *string `json:"last_name"`
	Displayname      string  `json:"displayname"`
	IsEmailConfirmed bool    `json:"-"`
	EmailConfirmHash string  `json:"-"`
}

type UserBasic struct {
	UserBasicModel
	Passwd      string `json:"-"`
}

type User42Model struct {
	User42Id      uint      `json:"-"`
	UserId		  uint		`json:"-"`
	AccessToken  *string    `json:"-"`
	RefreshToken *string    `json:"-"`
	ExpiresAt    *time.Time `json:"-"`
}

type User42 struct {
	User42Model
	Email       string `json:"-"`
	Fname       string `json:"-"`
	Lname       string `json:"-"`
	Displayname string `json:"-"`
	ImageBody   string `json:"-"`
}

type TokenHeader struct {
	UserId      uint   `json:"user_id"`
}

type Token struct {
	ServerPasswd string 	`json:"server_passwd,omitempty"`
	AccessToken  string 	`json:"access_token"`
}

func (user UserBasic) Validate() *errors.Error {
	if user.Email == "" || user.Passwd == "" {
		return errors.NoArgument.SetArgs("Email или пароль отсутствуют", "Email or password expected")
	}
	return nil
}

func (user *UserBasic) Sanitize() {
	user.Email = ""
}

func (user *UserBasic) ExtractFromUser42(user42 *User42) {
	user.User42Id = &user42.User42Id
	user.ImageBody = &user42.ImageBody
	user.Email = user42.Email
	user.Fname = &user42.Fname
	user.Lname = &user42.Lname
	user.Displayname = user42.Displayname
	user.IsEmailConfirmed = true
}
