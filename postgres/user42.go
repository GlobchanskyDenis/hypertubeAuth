package postgres

import (
	"HypertubeAuth/errors"
	"HypertubeAuth/model"
	"strconv"
	"strings"
)

func UserSet42(user *model.User42) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`INSERT INTO users_42_strategy (user_id, access_token, refresh_token,
		expires_at) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.UserId, user.AccessToken, user.RefreshToken, user.ExpiresAt)
	if err != nil {
		if strings.Contains(err.Error(), `users_42_strategy_pkey`) {
			return errors.ImpossibleToExecute.SetArgs("Такой пользователь уже существует в базе",
				"Such user already exists")
		}
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) != 1 {
		return errors.DatabaseExecutingError.SetArgs("добавлено "+strconv.Itoa(int(nbr64))+" пользователей",
			strconv.Itoa(int(nbr64))+" users was inserted")
	}
	return nil
}

func UserDelete42(user *model.User42) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`DELETE FROM users_42_strategy WHERE user_id = $1`)
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

func UserGet42ById(userId uint) (*model.User42, *errors.Error) {
	conn, Err := getConnection()
	if Err != nil {
		return nil, Err
	}
	stmt, err := conn.db.Prepare(`SELECT * FROM users_42_strategy WHERE user_id = $1`)
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
	var user = &model.User42{}
	if err := rows.Scan(&user.UserId, &user.AccessToken, &user.RefreshToken, &user.ExpiresAt); err != nil {
		return nil, errors.DatabaseScanError.SetOrigin(err)
	}
	return user, nil
}

func UserUpdate42(user *model.User42) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`UPDATE users_42_strategy
		SET access_token = $2, refresh_token = $3, expires_at = $4 WHERE user_id = $1;`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.UserId, user.AccessToken, user.RefreshToken, user.ExpiresAt)
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) == 0 {
		return errors.UserNotExist
	}
	if int(nbr64) > 1 {
		return errors.DatabaseExecutingError.SetArgs("обновлено "+strconv.Itoa(int(nbr64))+" пользователей",
			strconv.Itoa(int(nbr64))+" users was updated")
	}
	return nil
}
