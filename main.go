package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	dbmanager "github.com/antonvlasov/avito/dbmanager"
	netmanager "github.com/antonvlasov/avito/netmanager"
	//_ "github.com/mattn/go-sqlite3"
)

var q = make(chan subscription, 1000)
var stop bool = false

func updateAll() {
	items := dbmanager.GetItems()
	type record struct {
		id                 string
		oldPrice, newPrice float64
		rawurl             string
	}
	var toUpdate []record
	for id, oldPrice, rawurl, notLast := items.NextOldInfo(); notLast; id, oldPrice, rawurl, notLast = items.NextOldInfo() {
		toUpdate = append(toUpdate, record{id, oldPrice, 0, rawurl}) //can't modify a table while iterating through it
	}
	var wg sync.WaitGroup
	for i, rec := range toUpdate {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			toUpdate[i].newPrice, err = netmanager.GetPrice(rec.id)
			if err != nil {
				toUpdate[i].newPrice = toUpdate[i].oldPrice
				fmt.Println("request failed, not updating item")
			}
		}()
		wg.Wait()
	}
	for _, rec := range toUpdate {
		if rec.oldPrice != rec.newPrice {
			updateItem(rec.id, rec.newPrice, rec.rawurl)
		}
	}
}
func updateItem(id string, newPrice float64, rawurl string) {
	subs := dbmanager.GetSubscribers(id)
	for sub, notLast := subs.NextSubscriber(); notLast; sub, notLast = subs.NextSubscriber() {
		go netmanager.Notify(sub, rawurl, newPrice)
	}
	dbmanager.SetPrice(id, newPrice, rawurl)
}
func validate(rawurl string) bool {
	URL, err := url.Parse(rawurl)
	if err != nil {
		return false
	}
	if URL.Hostname() != "www.avito.ru" {
		return false
	}
	return true
}
func addSubscription(mail, url string) {
	if !validate(url) {
		return
		log.Fatal("invalid url, request ignored")
	}
	id, err := getItemID(url)
	if err != nil {
		log.Fatal("invalid url, request ignored")
	}
	oldPrice := dbmanager.GetPrice(id)
	newPrice, err := netmanager.GetPrice(id)
	if err != nil {
		newPrice = oldPrice
	}
	if oldPrice != newPrice {
		updateItem(id, newPrice, url)
	}
	dbmanager.AddSubscription(mail, id)
}
func EuqueueAddSubscription(mail, url string) {
	if len(q) < cap(q) {
		q <- subscription{mail, url}
	} else {
		fmt.Println("queue is full, request ignored")
	}
}

type subscription struct {
	mail, url string
}

//Run Runs the loop, updating all records every updateRate in seconds
//Multiply Runs should not be launched together
//If Runs is launched, it handles subscriptions queued by EuqueueAddSubscription
//To stop the loop call Stop()
func Run(updateRate int) {
	go func() {
		tick := time.Tick(time.Duration(updateRate) * time.Second)
		for {
			if stop {
				break
			}
			select {
			case <-tick:
				updateAll()
			case sub := <-q:
				addSubscription(sub.mail, sub.url)
			}
		}
	}()
}
func Stop() {
	stop = true
}
func getItemID(rawurl string) (string, error) {
	splitted := strings.Split(rawurl, "/")
	size := len(splitted)
	if size < 2 {
		return "", errors.New("incorrect url")
	}
	var s string
	if len(splitted[size-1]) != 0 {
		s = splitted[size-1]
	} else {
		s = splitted[size-2]
	}
	lastIndex := strings.LastIndex(s, "_")
	if lastIndex == -1 {
		return "", errors.New("incorrect url")
	}
	lastIndex += 1
	if lastIndex >= len(s) {
		return "", errors.New("incorrect url")
	}
	id := s[lastIndex:]
	return id, nil
}
func handlerHello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "use http://localhost:5000/api/subscribe?url=<item_url>&email=<mail> to subscribe")
}
func handlerSubscribeRequest(w http.ResponseWriter, req *http.Request) {
	method := req.URL.Path
	if method != "/api/subscribe" {
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid method")
		return
	}
	params := req.URL.Query()
	rawurl := params.Get("url")
	if rawurl == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid url")
		return
	}
	mail := params.Get("email")
	if mail == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid email")
		return
	}

	EuqueueAddSubscription(mail, rawurl)
	fmt.Fprint(w, "you have been subscribed")
}

func main() {
	dbmanager.Connect("./resources/")
	Run(600)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/subscribe", handlerSubscribeRequest)
	mux.HandleFunc("/", handlerHello)
	http.ListenAndServe(":5000", mux)
}
