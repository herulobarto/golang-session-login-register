package models

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"github.com/herulobarto/go-auth/config"
	"github.com/herulobarto/go-auth/entities"
)

type UserModel struct {
	db *sql.DB
}

func NewUserModel() *UserModel {
	conn, err := config.DBConn()
	if err != nil {
		panic(err)
	}

	return &UserModel{
		db: conn,
	}
}

func (u UserModel) Where(user *entities.User, fieldName, fieldValue string) error {
	query := "SELECT id, nama_lengkap, email, username, password FROM users WHERE " + fieldName + " = ? LIMIT 1"
	row := u.db.QueryRow(query, fieldValue)

	err := row.Scan(&user.Id, &user.NamaLengkap, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no user found")
		}
		return err
	}

	return nil
}

func (u UserModel) Create(user entities.User) (int64, error) {
	result, err := u.db.Exec("INSERT INTO users (nama_lengkap, email, username, password) VALUES (?, ?, ?, ?)",
		user.NamaLengkap, user.Email, user.Username, user.Password)

	if err != nil {
		return 0, err
	}

	lastInsertId, _ := result.LastInsertId()

	return lastInsertId, nil
}
