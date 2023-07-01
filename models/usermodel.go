package models

import (
	"database/sql"

	"github.com/taufiqkba/go_auth/config"
	"github.com/taufiqkba/go_auth/entities"
)

type UserModel struct {
	db *sql.DB
}

func NewUserModel() *UserModel {
	conn, err := config.DBConn()
	if err != nil {
		panic(err)
	}
	return &UserModel{conn}
}

func (u UserModel) Where(user *entities.User, fieldName, fieldValue string) error {
	row, err := u.db.Query("SELECT * FROM users WHERE "+fieldName+" = ? limit 1 ", fieldValue)
	if err != nil {
		return err
	}

	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {

		}
	}(row)

	for row.Next() {
		row.Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Password)
	}
	return nil
}

func (u UserModel) Create(user entities.User) (int64, error) {
	result, err := u.db.Exec("INSERT INTO users (name, email, username, password) values (?,?,?,?)", user.Name, user.Email, user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	lastInsertId, _ := result.LastInsertId()
	return lastInsertId, nil
}
