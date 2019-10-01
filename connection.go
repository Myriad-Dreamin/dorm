package dorm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
}


func Open() (*DB, error) {
	db, err := sql.Open("mysql", "root:111111@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}




