package postgres

import (
	"HypertubeAuth/errors"
	"HypertubeAuth/model"
	"strconv"
	"strings"
)

func UserSetBasic(user *model.UserBasic) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`INSERT INTO users (email, encryptedpass, first_name, last_name,
		displayname, email_confirm_hash) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	if err = stmt.QueryRow(user.Email, user.EncryptedPass, user.Fname, user.Lname, user.Displayname,
		user.EmailConfirmHash).Scan(&user.UserId); err != nil {
		if strings.Contains(err.Error(), `users_email_key`) {
			return errors.ImpossibleToExecute.SetArgs("Эта почта уже закреплена за другим пользователем",
				"This email is already assigned to another user")
		}
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	return nil
}

func UserDeleteBasic(user *model.UserBasic) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`DELETE FROM users WHERE user_id = $1`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	result, err := stmt.Exec(user.UserId)
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) == 0 {
		return errors.ImpossibleToExecute.SetArgs("Пользователь не найден", "User not found")
	}
	if int(nbr64) > 1 {
		return errors.DatabaseExecutingError.SetArgs("удалено "+strconv.Itoa(int(nbr64))+" пользователя",
			strconv.Itoa(int(nbr64))+" users was deleted")
	}
	return nil
}

func UserGetBasicById(userId uint) (*model.UserBasic, *errors.Error) {
	conn, Err := getConnection()
	if Err != nil {
		return nil, Err
	}
	stmt, err := conn.db.Prepare(`SELECT * FROM users WHERE user_id = $1`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.UserNotExist
	}
	var user = &model.UserBasic{}
	if err := rows.Scan(&user.UserId, &user.ImageBody, &user.Email, &user.EncryptedPass, &user.Fname,
		&user.Lname, &user.Displayname, &user.IsEmailConfirmed, &user.EmailConfirmHash); err != nil {
		return nil, errors.DatabaseScanError.SetOrigin(err)
	}
	return user, nil
}

func UserGetBasicByEmail(email string) (*model.UserBasic, *errors.Error) {
	conn, Err := getConnection()
	if Err != nil {
		return nil, Err
	}
	stmt, err := conn.db.Prepare(`SELECT * FROM users WHERE email = $1`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(email)
	if err != nil {
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.UserNotExist
	}
	var user = &model.UserBasic{}
	if err := rows.Scan(&user.UserId, &user.ImageBody, &user.Email, &user.EncryptedPass, &user.Fname,
		&user.Lname, &user.Displayname, &user.IsEmailConfirmed, &user.EmailConfirmHash); err != nil {
		return nil, errors.DatabaseScanError.SetOrigin(err)
	}
	return user, nil
}

func UserUpdateBasic(user *model.UserBasic) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`UPDATE users SET is_email_confirmed = true WHERE user_id = $1`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	result, err := stmt.Exec(user.UserId)
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) == 0 {
		return errors.ImpossibleToExecute.SetArgs("Пользователь не найден", "User not found")
	}
	if int(nbr64) > 1 {
		return errors.DatabaseExecutingError.SetArgs("изменено "+strconv.Itoa(int(nbr64))+" пользователя",
			strconv.Itoa(int(nbr64))+" users was updated")
	}
	return nil
}
