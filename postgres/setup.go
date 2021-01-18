package postgres

import (
	"HypertubeAuth/errors"
	"database/sql"
	_ "github.com/lib/pq"
)

type Connection struct {
	db *sql.DB
}

var gConnection *Connection

func Init() *errors.Error {
	if cfg == nil {
		return errors.NotConfiguredPackage.SetArgs("postgres", "postgres")
	}

	var conn Connection
	var err error
	dsn := "user=" + cfg.User + " password=" + cfg.Passwd + " dbname=" + cfg.Database + " host=" + cfg.Host + " sslmode=disable"
	if conn.db, err = sql.Open(cfg.Type, dsn); err != nil {
		return errors.DatabaseError.SetArgs("Не смог установить соединение", "Connection fail").SetOrigin(err)
	}
	if err = conn.db.Ping(); err != nil {
		return errors.DatabaseError.SetArgs("Ping к БД вернул ошибку", "Ping DB returned error").SetOrigin(err)
	}
	gConnection = &conn
	return nil
}

func Close() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if err := conn.db.Close(); err != nil {
		return errors.DatabaseError.SetArgs("Не смог закрыть соединение", "Connection close failed").SetOrigin(err)
	}
	return nil
}

func DropAllTables() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if _, err := conn.db.Exec("DROP TABLE IF EXISTS images"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	if _, err := conn.db.Exec("DROP TABLE IF EXISTS users_42_strategy"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	if _, err := conn.db.Exec("DROP TABLE IF EXISTS users"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	return nil
}

func CreateUsersTable() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if _, err := conn.db.Exec("CREATE TABLE users(user_id SERIAL PRIMARY KEY, " +
		"image_body VARCHAR, " +
		"email VARCHAR CONSTRAINT users_email_key UNIQUE NOT NULL, " +
		"encryptedPass VARCHAR(35) NOT NULL, " +
		"first_name VARCHAR NOT NULL, " +
		"last_name VARCHAR NOT NULL, " +
		"displayname VARCHAR NOT NULL, " +
		"is_email_confirmed BOOL NOT NULL DEFAULT false, " +
		"email_confirm_hash VARCHAR NOT NULL)"); err != nil {
		return errors.DatabaseError.SetArgs("11", "11").SetOrigin(err)
	}
	return nil
}

func CreateUsers42StrategyTable() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if _, err := conn.db.Exec("CREATE TABLE users_42_strategy(user_id INTEGER PRIMARY KEY, " +
		"access_token VARCHAR, " +
		"refresh_token VARCHAR, " +
		"expires_at TIMESTAMP)"); err != nil {
		return errors.DatabaseError.SetArgs("2", "2").SetOrigin(err)
	}
	return nil
}
