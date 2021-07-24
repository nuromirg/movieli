package service

import (
	"Movieli/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB(dbName string) *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/" + dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_,err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	if err != nil {
		panic(err)
	}
	_,err = db.Exec("USE "+ dbName + ";")
	if err != nil {
		panic(err)
	}
	_,err = db.Exec("CREATE TABLE IF NOT EXISTS movieli (id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, poster VARCHAR(256), title VARCHAR(128), year INTEGER, director VARCHAR(128));")
	if err != nil {
		panic(err)
	}
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/" + dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return db
}

func DBConnect() (db *sql.DB) {
	db, err := sql.Open("mysql", "root@/" + config.DBNAME)
	if err != nil {
		panic(err.Error())
	}
	return db
}

