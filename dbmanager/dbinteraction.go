package dbmanager

import (
	"database/sql"
)

func SetPrice(id string, price float64, rawurl string) {
	if db == nil {
		panic("db is nil")
		return
	}
	db.Exec("INSERT OR IGNORE INTO Items VALUES($1,$2,$3); UPDATE Items SET price = $4 WHERE id = $5;", id, price, rawurl, price, id)
}

//AddSubscription implies there is already a record in Items table
func AddSubscription(mail string, id string) {
	db.Exec("INSERT OR IGNORE INTO Subscribers VALUES($1)", mail)
	db.Exec("INSERT OR IGNORE INTO Subscriptions VALUES($1,$2)", mail, id)

}
func GetPrice(id string) float64 {
	if db == nil {
		panic("db is nil")
		return 0
	}
	row := db.QueryRow("SELECT price FROM Items WHERE id = $1", id)
	var price float64
	err := row.Scan(&price)
	switch {
	case err == sql.ErrNoRows:
		return -1
	default:
		return price
	}
}

type Items struct {
	rows *sql.Rows
}

func GetItems() *Items {
	if db == nil {
		panic("db is nil")
		return nil
	}
	res := &Items{}
	res.rows, _ = db.Query("SELECT * FROM Items;")
	return res
}
func (items *Items) NextOldInfo() (id string, price float64, rawurl string, notLast bool) {
	if !items.rows.Next() {
		return
	}
	items.rows.Scan(&id, &price, &rawurl)
	notLast = true
	return
}

type Subscribers struct {
	rows *sql.Rows
}

func GetSubscribers(id string) *Subscribers {
	if db == nil {
		panic("db is nil")
		return nil
	}
	res := &Subscribers{}
	res.rows, _ = db.Query("SELECT mail FROM Subscriptions WHERE id = $1;", id)
	return res
}

func (subscibers *Subscribers) NextSubscriber() (mail string, notLast bool) {
	if !subscibers.rows.Next() {
		return
	} else {
		subscibers.rows.Scan(&mail)
		notLast = true
		return
	}
}
