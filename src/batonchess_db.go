package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	sqliteVersion = "sqlite3"
	userDb        = "./db/users.db"
)

// --- SETUP

func setupDb() error {
	err := createUserDb()
	if err != nil {
		return err
	}

	return nil
}

func createUserDb() error {
	return execNonQuery(userDb, `
		CREATE TABLE users (
		Id TEXT PRIMARY KEY,
		Name TEXT NOT NULL);`,
	)
}

// --- DB

func execNonQuery(dbLocation string, queryString string, queryArgs ...any) error {
	db, err := sql.Open(sqliteVersion, dbLocation)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	stmt, err := tx.Prepare(queryString)
	if err != nil {
		println(err.Error())
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(queryArgs...)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func execQuery(dbLocation string, queryString string, queryArgs ...any) (int, error) {
	return 0, nil
}

// --- USER

func InsertUser(user *User) error {
	return execNonQuery(
		userDb,
		`INSERT INTO users(Id, Name) VALUES(?, ?)`,
		user.Id, user.Name,
	)
}

func UpdateUserName(updateInfo *UserNameUpdateRequest) error {
	return execNonQuery(
		userDb,
		`UPDATE users SET Name = ? WHERE Id = ?`,
		updateInfo.NewUsername, updateInfo.UserId,
	)
}
