package dbmanager

import (
	"errors"
	"log"
)

type Item struct {
	Id           string
	Price        float64
	Url          string
	ChangeStatus int
}
type Subscription struct {
	Mail  string
	Id    string
	IsNew bool
}
type ToVerify struct {
	Link, Mail string
}

func ReadAllItems() (items []Item, err error) {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	rows, err := db.Query("SELECT * FROM Items")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item Item
		rows.Scan(&item.Id, &item.Price, &item.Url, &item.ChangeStatus)
		items = append(items, item)
	}
	return items, nil
}
func IsVerified(mail string) (isVerified bool, err error) {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	rows, err := db.Query("SELECT mail FROM Subscribers WHERE mail = $1", mail)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}
	return false, nil
}
func AddMail(mail string) error {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	_, err := db.Exec("INSERT OR IGNORE INTO Subscribers VALUES($1)", mail)
	return err
}

//Verify adds mail to Subscribers automatically
func Verify(link string) (verified bool, err error) {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	rows, err := db.Query("SELECT * FROM ToVerify WHERE link = $1", link)
	if err != nil {
		return
	}
	defer rows.Close()
	var mail string
	for rows.Next() {
		rows.Scan(&link, &mail)
		verified = true
		break
	}
	rows.Close()
	if !verified {
		return false, nil
	}
	err = AddMail(mail)
	if err != nil {
		return false, err
	}
	_, err = db.Exec("DELETE FROM ToVerify WHERE link = $1", link)
	if err != nil {
		return true, nil
	}
	return true, nil

}
func AddVerificationLink(link, mail string) error {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	_, err := db.Exec("DELETE FROM ToVerify WHERE mail = $1", mail)
	rows, err := db.Query("SELECT link FROM ToVerify WHERE link = $1", link)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		return errors.New("such link already exists")
	}
	_, err = db.Exec("INSERT INTO ToVerify VALUES($1,$2)", link, mail)
	if err != nil {
		return err
	}
	return nil
}
func UpdateOrInsertItem(item Item) error {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	_, err := db.Exec("INSERT OR IGNORE INTO Items VALUES($1,$2,$3,$4)", item.Id, item.Price, item.Url, item.ChangeStatus)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE Items SET id = $1, price = $2, url = $3, change_status = $4  WHERE id = $5",
		item.Id, item.Price, item.Url, item.ChangeStatus, item.Id)
	if err != nil {
		return err
	}
	return nil
}
func ReadSubscriptionsForItem(id string) (subscriptions []Subscription, err error) {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	rows, err := db.Query("SELECT * FROM Subscriptions WHERE id = $1", id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sbscr Subscription
		rows.Scan(&sbscr.Mail, &sbscr.Id, &sbscr.IsNew)
		subscriptions = append(subscriptions, sbscr)
	}
	return subscriptions, nil
}
func AddSubscription(sbscr Subscription) error {
	if db == nil {
		log.Fatal("db pointer is nil")
	}
	rows, err := db.Query("SELECT mail FROM Subscriptions WHERE (mail = $1 AND id = $2)", sbscr.Mail, sbscr.Id)
	if err != nil {
		return err
	}
	for rows.Next() {
		rows.Close()
		return nil
	}
	_, err = db.Exec("INSERT INTO Subscriptions VALUES($1,$2,$3)", sbscr.Mail, sbscr.Id, sbscr.IsNew)
	if err != nil {
		return err
	}
	return nil
}
func MakeSubscriptionOld(subscription Subscription) error {
	if db == nil {
		panic("db is nil")
	}
	mail := subscription.Mail
	id := subscription.Id
	_, err := db.Exec("UPDATE Subscriptions SET mail = $1, id = $2, is_new = $3 WHERE mail = $4 AND id = $5", mail, id, false, mail, id)
	if err != nil {
		return err
	}
	return nil
}

// func SetPrice(id string, price float64, rawurl string) {
// 	if db == nil {
// 		panic("db is nil")
// 		return
// 	}
// 	db.Exec("INSERT OR IGNORE INTO Items VALUES($1,$2,$3); UPDATE Items SET price = $4 WHERE id = $5;", id, price, rawurl, price, id)
// }

// //AddSubscription implies there is already a record in Items table
// func AddSubscription(mail string, id string) {
// 	db.Exec("INSERT OR IGNORE INTO Subscribers VALUES($1)", mail)
// 	db.Exec("INSERT OR IGNORE INTO Subscriptions VALUES($1,$2)", mail, id)

// }
// func GetPrice(id string) float64 {
// 	if db == nil {
// 		panic("db is nil")
// 		return 0
// 	}
// 	row := db.QueryRow("SELECT price FROM Items WHERE id = $1", id)
// 	var price float64
// 	err := row.Scan(&price)
// 	switch {
// 	case err == sql.ErrNoRows:
// 		return -1
// 	default:
// 		return price
// 	}
// }

// type Items struct {
// 	rows *sql.Rows
// }

// func GetItems() *Items {
// 	if db == nil {
// 		panic("db is nil")
// 		return nil
// 	}
// 	res := &Items{}
// 	res.rows, _ = db.Query("SELECT * FROM Items;")
// 	return res
// }
// func (items *Items) NextOldInfo() (id string, price float64, rawurl string, notLast bool) {
// 	if !items.rows.Next() {
// 		return
// 	}
// 	items.rows.Scan(&id, &price, &rawurl)
// 	notLast = true
// 	return
// }

// type Subscribers struct {
// 	rows *sql.Rows
// }

// func GetSubscribers(id string) *Subscribers {
// 	if db == nil {
// 		panic("db is nil")
// 		return nil
// 	}
// 	res := &Subscribers{}
// 	res.rows, _ = db.Query("SELECT mail FROM Subscriptions WHERE id = $1;", id)
// 	return res
// }

// func (subscibers *Subscribers) NextSubscriber() (mail string, notLast bool) {
// 	if !subscibers.rows.Next() {
// 		return
// 	} else {
// 		subscibers.rows.Scan(&mail)
// 		notLast = true
// 		return
// 	}
// }
