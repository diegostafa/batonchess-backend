package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func GetUser(u *UserId) (*User, error) {
	var (
		user User
		err  error
	)
	err = queryOne(
		func(row *sql.Row) error {
			return row.Scan(&user.Id, &user.Name)
		}, `
		SELECT *
		FROM users
		WHERE u_id = ?`,
		u.Id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(user *User) error {
	return queryNone(`
		INSERT INTO users(u_id, u_name)
		VALUES (?, ?)`,
		user.Id, user.Name)
}

func UpdateUserName(updateInfo *UpdateUsernameRequest) error {
	return queryNone(`
		UPDATE users
		SET u_name = ?
		WHERE u_id = ?`,
		updateInfo.NewName, updateInfo.Id,
	)
}
