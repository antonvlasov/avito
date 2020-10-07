package dbmanager

import (
	"testing"
)

func init() {
	path := "../testdata/"
	Connect(path)
}

func clearDB() {
	if db == nil {
		panic("db is nil")
	}
	db.Exec("DELETE FROM Subscriptions; DELETE FROM Items; DELETE FROM Subscribers;")
}
func fillDB() {
	if db == nil {
		panic("db is nil")
	}
	db.Exec(`INSERT INTO Items VALUES(1,300),(2,700),(3,500),(4,600);
			 INSERT INTO Subscribers VALUES('mail1'),('mail2'),'(mail3');
			 INSER INTO Subscriptions VALUES('mail1',1),('mail2',2),('mail2',1),('mail3',2),('mail3',3),('mail3',4)`)
}
func TestSetPrice(t *testing.T) {
	clearDB()
	CheckResults := func(id int, p float64) {
		rows, _ := db.Query("SELECT * FROM Items WHERE id = $1", id, p)
		for rows.Next() {
			var priceRes float64
			var idRes int
			rows.Scan(&idRes, &priceRes)
			if priceRes != p {
				t.Errorf("SetPrice failed, expected price = %v, got %v", p, priceRes)
			}
		}
	}
	var it1 = 1
	var it2 = 2
	var p1 = 100.0
	var p2 = 200.0
	SetPrice(it1, p1)
	CheckResults(it1, p1)
	SetPrice(it2, p2)
	CheckResults(it1, p1)
	CheckResults(it2, p2)
	SetPrice(it1, p2)
	CheckResults(it1, p2)
	CheckResults(it2, p2)
}

func TestAddSubscription(t *testing.T) {
	fillDB()
	// CheckResults := func(id int, p float64) {
	// 	rows, _ := db.Query("SELECT * FROM Items WHERE id = $1", id, p)
	// 	for rows.Next() {
	// 		var priceRes float64
	// 		var idRes int
	// 		rows.Scan(&idRes, &priceRes)
	// 		if priceRes != p {
	// 			t.Errorf("SetPrice failed, expected price = %v, got %v", p, priceRes)
	// 		}
	// 	}
	// }
}
