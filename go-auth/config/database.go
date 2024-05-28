package config

import (
	"database/sql"
)

func DBConn() (db *sql.DB, err error) {

	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "go-auth"

	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	return

}
