package batonchess

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// --- USER

func GetUser(u *UserId) (*User, error) {
	var (
		user User
		err  error
	)
	err = queryOne(
		func(row *sql.Row) error {
			return row.Scan(&user.Id, &user.Name)
		}, `
		SELECT
			*
		FROM
			users
		WHERE
			u_id = ?`,
		u.Id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(user *User) error {
	return queryNone(`
		INSERT INTO
			users(u_id, u_name)
		VALUES
			(?, ?)`,
		user.Id, user.Name)
}

func UpdateUserName(updateInfo *UserNameUpdateRequest) error {
	return queryNone(`
		UPDATE
			users
		SET
			u_name= ?
		WHERE
			u_id = ?`,
		updateInfo.NewUsername, updateInfo.Id,
	)
}
