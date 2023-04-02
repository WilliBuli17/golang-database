package golang_database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "user1:password@tcp(localhost:3306)/golang_database?parseTime=true")

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(25)
	db.SetMaxIdleConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
