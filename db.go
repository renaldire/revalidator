package Validator

import "database/sql"

var db *sql.DB

func getDb() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open(DbDriver, ConnectionString)
		if err != nil {
			panic(err.Error())
		}
	}
	return db
}

