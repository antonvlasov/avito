package dbmanager

import (
	"fmt"
	"testing"
)

const path = "../testdata/"

func TestConnectDisconnect(t *testing.T) {
	Connect(path)
	if db == nil {
		Close()
		t.Errorf("ConnectDisconnect failed db pointer is nil")
	}
	Close()
}
func TestItems(t *testing.T) {
	Connect(path)
	Cleardb()
	items := []Item{
		{"1867132320", 400, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya_kanonicheskaya_1867132320", 0},
		{"1897268986", 1500, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya._kanonicheskie_pisaniya_1897268986", 0},
		{"1897268986", 1000, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya._kanonicheskie_pisaniya_1897268986", 1},
		Item{},
	}
	for i, item := range items {
		err := UpdateOrInsertItem(item)
		if err != nil {
			t.Errorf("UpdateOrInsertItem failed on item %v", i)
		}
	}
	correctItems := []Item{
		{"1867132320", 400, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya_kanonicheskaya_1867132320", 0},
		{"1897268986", 1000, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya._kanonicheskie_pisaniya_1897268986", 1},
		Item{},
	}
	readItems, err := ReadAllItems()
	if err != nil {
		t.Errorf("ReadAllItems failed with err %v", err)
	}
	if len(correctItems) != len(readItems) {
		t.Errorf("ReadAllItems failed, different len")
	}
	for i := range correctItems {
		if correctItems[i] != readItems[i] {
			t.Errorf("ReadAllItems failed, different items: expected %v got %v", correctItems[i], readItems[i])
		}
	}
	fmt.Println("passed")
	Close()
}
func TestSubscribers(t *testing.T) {
	Connect(path)
	Cleardb()
	subs := []string{"mail1", "mail2", "", "mail1"}
	ans := []bool{false, false, false, true}
	for i := range subs {
		isV, err := IsVerified(subs[i])
		if err != nil {
			t.Errorf("IsVerified failed, got error %v", err)
		}
		if isV != ans[i] {
			t.Errorf("IsVerified failed: expected %v got %v", ans[i], isV)
		}
		err = AddMail(subs[i])
		if err != nil {
			t.Errorf("AddMail failed, got error %v", err)
		}
		isV, err = IsVerified(subs[i])
		if err != nil {
			t.Errorf("IsVerified failed, got error %v", err)
		}
		if !isV {
			t.Errorf("IsVerified failed: expected %v got %v", true, isV)
		}
	}
	Close()
}

func TestVerification(t *testing.T) {
	Connect(path)
	Cleardb()
	toVerify := []ToVerify{
		{"link1", "mail1"},
		{"link2", "mail2"},
		{"link2", "mail3"},
		{"link3", "mail1"},
	}
	links := []string{"link1", "link2", "link3", "link4"}
	verifyResults := []bool{false, true, true, false}
	mails := []string{"mail1", "mail2", "mail3"}
	mailCheckResults := [][]bool{{false, false, false}, {false, true, false}, {true, true, false}, {true, true, false}, {true, true, false}}
	for i := range toVerify {
		err := AddVerificationLink(toVerify[i].Link, toVerify[i].Mail)
		if err != nil {
			if err.Error() != "such link already exists" {
				t.Errorf("AddVerificationLink failed: error %v", err)
			}
		}
	}
	for i := range links {
		ver, err := Verify(links[i])
		if err != nil {
			t.Errorf("Verify failed: error %v", err)
		}
		if ver != verifyResults[i] {
			t.Errorf("Verify failed on %v: expected %v", links[i], verifyResults[i])
		}
		for j := range mails {
			isV, err := IsVerified(mails[j])
			if err != nil {
				t.Errorf("IsVerified failed: error %v", err)
			}
			if isV != mailCheckResults[i][j] {
				t.Errorf("AddVerificationLink failed on iteration %v: expected verification of %v to be %v", i, mails[j], mailCheckResults[i][j])
			}
		}
	}
	Close()
}
func TestSubscription(t *testing.T) {
	Connect(path)
	Cleardb()
	items := []Item{
		{"1867132320", 400, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya_kanonicheskaya_1867132320", 0},
		{"1897268986", 1500, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya._kanonicheskie_pisaniya_1897268986", 0},
		{"1897268986", 1000, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya._kanonicheskie_pisaniya_1897268986", 1},
		Item{},
	}
	for i, item := range items {
		err := UpdateOrInsertItem(item)
		if err != nil {
			t.Errorf("UpdateOrInsertItem failed on item %v", i)
		}
	}
	subs := []string{"mail1", "mail2", "", "mail1"}
	for i := range subs {
		err := AddMail(subs[i])
		if err != nil {
			t.Errorf("AddMail failed, got error %v", err)
		}
	}
	subscriptions := []Subscription{
		{"mail1", "1867132320", true},
		{"mail1", "1867132320", true},
		{"mail2", "1897268986", false},
		{"mail2", "1867132320", true},
		{"", "", false},
		{"not an email", "not an item", false},
	}
	ids := []string{"1867132320", "1897268986", "", "not in table"}
	correctSubscriptions := [][]Subscription{
		{{"mail1", "1867132320", true}, {"mail2", "1867132320", true}},
		{{"mail2", "1897268986", false}},
		{{"", "", false}},
		{},
	}
	for i := range subscriptions {
		err := AddSubscription(subscriptions[i])
		if err != nil {
			if i != 5 {
				t.Errorf("AddSubscription failed, got error %v", err)
			}
		}
	}
	for i := range items {
		readSubscriptions, err := ReadSubscriptionsForItem(ids[i])
		if err != nil {
			t.Errorf("ReadSubscriptionsForItem failed, gor error %v", err)
		}
		if len(readSubscriptions) != len(correctSubscriptions[i]) {
			t.Errorf("ReadSubscriptionsForItem failed, different len")
		}
		for j := range correctSubscriptions[i] {
			if readSubscriptions[j] != correctSubscriptions[i][j] {
				t.Errorf("ReadSubscriptionsForItem failed, got wrong subscription")
			}
		}
	}
	for i := range subscriptions {
		err := MakeSubscriptionOld(subscriptions[i])
		if err != nil {
			t.Errorf("MakeSubscriptionOld failed, got error %v", err)
		}
	}
	for i := range items {
		readSubscriptions, err := ReadSubscriptionsForItem(ids[i])
		if err != nil {
			t.Errorf("ReadSubscriptionsForItem failed, gor error %v", err)
		}
		if len(readSubscriptions) != len(correctSubscriptions[i]) {
			t.Errorf("ReadSubscriptionsForItem failed, different len")
		}
		for j := range correctSubscriptions[i] {
			if readSubscriptions[j].IsNew != false {
				t.Errorf("MakeSubscriptionOld failed, subscription %v %v is not old", readSubscriptions[j].Mail, readSubscriptions[j].Id)
			}
		}
	}
}

func Cleardb() {
	if db == nil {
		panic("db is nil")
	}
	db.Exec("DELETE FROM Subscriptions; DELETE FROM Items; DELETE FROM Subscribers; DELETE FROM ToVerify")
}
