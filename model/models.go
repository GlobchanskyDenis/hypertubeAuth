package model

import (
	"HypertubeAuth/errors"
	"time"
)

type UserBasicModel struct {
	UserId           uint    `json:"user_id"`
	ImageBody        *string `json:"image_body"`
	Email            string  `json:"email"`
	EncryptedPass    string  `json:"-"`
	Fname            string  `json:"first_name"`
	Lname            string  `json:"last_name"`
	Displayname      string  `json:"displayname"`
	IsEmailConfirmed bool    `json:"-"`
	EmailConfirmHash string  `json:"-"`
}

type UserBasic struct {
	UserBasicModel
	Passwd      string `json:"-"`
}

type User42Model struct {
	UserId       uint       `json:"user_id"`
	AccessToken  *string    `json:"-"`
	RefreshToken *string    `json:"-"`
	ExpiresAt    *time.Time `json:"-"`
}

type User42 struct {
	User42Model
	Email       string `json:"email"`
	Fname       string `json:"first_name"`
	Lname       string `json:"last_name"`
	Displayname string `json:"displayname"`
	ImageBody   string `json:"image_body"`
}

type User struct {
	UserId      uint    `json:"user_id"`
	Email       string  `json:"email"`
	Fname       string  `json:"first_name"`
	Lname       string  `json:"last_name"`
	Displayname string  `json:"displayname"`
	ImageBody   *string `json:"image_body"`
}

type TokenHeader struct {
	UserId      uint   `json:"user_id"`
	AccountType string `json:"account_type"`
}

type Token struct {
	ServerPasswd string `json:"server_passwd,omitempty"`
	AccessToken  string `json:"access_token"`
	Profile      *User  `json:"profile,omitempty"`
}

func (user UserBasic) Validate() *errors.Error {
	if user.Email == "" || user.Passwd == "" {
		return errors.NoArgument.SetArgs("Email или пароль отсутствуют", "Email or password expected")
	}
	return nil
}

func (user UserBasic) TransformToUser() User {
	return User{
		UserId: user.UserId,
		Email: user.Email,
		Fname: user.Fname,
		Lname: user.Lname,
		Displayname: user.Displayname,
		ImageBody: user.ImageBody,
	}
}

func (user User42) TransformToUser() User {
	return User{
		UserId: user.UserId,
		Email: user.Email,
		Fname: user.Fname,
		Lname: user.Lname,
		Displayname: user.Displayname,
		ImageBody: &user.ImageBody,
	}
}
