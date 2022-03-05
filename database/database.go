package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	connectionData := "golang:golang@/devbook?charset=utf8&parseTime=True&loc=Local"
	database, error := sql.Open("mysql", connectionData)
	if error != nil {
		return nil, error
	}
	if error = database.Ping(); error != nil {
		return nil, error
	}
	return database, nil
}
