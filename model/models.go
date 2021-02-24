package model

import (
	"HypertubeAuth/controller/validator"
	"HypertubeAuth/errors"
	"time"
)

type UserBasicModel struct {
	UserId           uint    `json:"userId"`
	User42Id         *uint   `json:"-"`
	ImageBody        *string `json:"imageBody"`
	Email            string  `json:"email"`
	EncryptedPass    *string `json:"-"`
	Fname            *string `json:"firstName"`
	Lname            *string `json:"lastName"`
	Username         string  `json:"username"`
	IsEmailConfirmed bool    `json:"-"`
	NewEmail		 *string `json:"-"`
}

type UserBasic struct {
	UserBasicModel
	Passwd string `json:"-"`
}

type User42Model struct {
	User42Id     uint       `json:"-"`
	UserId       uint       `json:"-"`
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

type AccessTokenHeader struct {
	UserId uint `json:"userId"`
}

type RepairTokenHeader struct {
	UserId uint `json:"userId"`
}

type EmailTokenHeader struct {
	UserId   uint   `json:"userId"`
	NewEmail string `json:"newEmail"`
}

type Token struct {
	ServerPasswd string `json:"serverPasswd,omitempty"`
	AccessToken  string `json:"accessToken"`
}

func (user UserBasic) Validate() *errors.Error {
	if user.Email == "" || user.Passwd == "" {
		return errors.NoArgument.SetArgs("Email или пароль отсутствуют", "Email or password expected")
	}
	if Err := validator.ValidateEmail(user.Email); Err != nil {
		return Err
	}
	if Err := validator.ValidatePassword(user.Passwd); Err != nil {
		return Err
	}
	if Err := validator.ValidateName(user.Username); Err != nil {
		return Err
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
	user.Username = user42.Displayname
	user.IsEmailConfirmed = true
}
