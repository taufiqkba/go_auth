package config

import "database/sql"

func DB() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPassword := "root"
	dbName := "go_auth"

	db, err = sql.Open(dbDriver, dbUser+":"+dbPassword+"@/"+dbName)
	return
}
