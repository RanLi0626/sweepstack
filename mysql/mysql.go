package mysql

import (
	"database/sql"
	"log"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// GetConn get the connection for mysql
func GetConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@/sweepstake")
	if err != nil {
		log.Println("connect mysql error", err)
		return nil, err
	}
	return db, nil
}
