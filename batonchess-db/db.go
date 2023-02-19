package batonchessDb

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	sqliteVersion = "sqlite3"
	batonchessDb  = "./db/batonchess.db"
)

// --- DB

func queryNone(queryString string, queryArgs ...any) error {
	db, err := sql.Open(sqliteVersion, batonchessDb)
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

func queryOne(scanner func(*sql.Row) error, queryString string, queryArgs ...any) error {
	db, err := sql.Open(sqliteVersion, batonchessDb)
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

	row := stmt.QueryRow(queryArgs...)
	scanned := scanner(row)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return scanned
}

func queryMany(scanner func(*sql.Rows) error, queryString string, queryArgs ...any) error {
	db, err := sql.Open(sqliteVersion, batonchessDb)
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

	rows, err := stmt.Query(queryArgs...)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}

	scanned := scanner(rows)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return scanned
}
