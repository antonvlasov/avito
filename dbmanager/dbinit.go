package dbmanager

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func base() {
	fmt.Println("init")
}

func Connect(path string) {
	var err error
	db, err = sql.Open("sqlite3", path+"Avito.db")

	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	PRAGMA foreign_keys = ON;
	PRAGMA journal_mode = WAL;

	CREATE TABLE IF NOT EXISTS Items(
	  id TEXT PRIMARY KEY NOT NULL,
	  price REAL NOT NULL,
	  url TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS Subscribers(
	  mail TEXT PRIMARY KEY NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS Subscriptions(
		mail TEXT NOT NULL,
		id TEXT NOT NULL,
		FOREIGN KEY(mail) REFERENCES Subscribers(mail),
		FOREIGN KEY(id) REFERENCES Items(id),
		UNIQUE(mail,id)
	)
	`)
	if err != nil {
		panic(err)
	}
}

func Close() {
	db.Close()
}

func ClearDB() {
	if db == nil {
		panic("db is nil")
	}
	db.Exec("DELETE FROM Subscriptions; DELETE FROM Items; DELETE FROM Subscribers;")
}
