package dbmanager

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Connect(path string) {
	var err error
	db, err = sql.Open("sqlite3", path+"Avito.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`
	PRAGMA foreign_keys = ON;
	PRAGMA journal_mode = WAL;

	CREATE TABLE IF NOT EXISTS Items(
	  id TEXT PRIMARY KEY NOT NULL,
	  price REAL NOT NULL,
	  url TEXT NOT NULL,
	  change_status INT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS Subscribers(
	  mail TEXT PRIMARY KEY NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS Subscriptions(
		mail TEXT NOT NULL,
		id TEXT NOT NULL,
		is_new BOOL NOT NULL,
		FOREIGN KEY(mail) REFERENCES Subscribers(mail),
		FOREIGN KEY(id) REFERENCES Items(id),
		UNIQUE(mail,id)
	);
	
	CREATE TABLE IF NOT EXISTS ToVerify(
		link TEXT PRIMARY KEY NOT NULL,
		mail TEXT NOT NULL,
		UNIQUE(mail)
	)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	db.Close()
}

func ClearDB() {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	db.Exec("DELETE FROM Subscriptions; DELETE FROM Items; DELETE FROM Subscribers;")
}
