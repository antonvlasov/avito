package httpserver

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/mail"
	"strings"

	"github.com/antonvlasov/avito/netmanager"

	"github.com/antonvlasov/avito/dbmanager"
)

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/subscribe", handleSubscribeRequest)
	mux.HandleFunc("/", handleFrontPage)
	mux.HandleFunc("/verify", handleVerification)
	http.ListenAndServe(":5000", mux)
}
func CheckEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
func GetItemID(rawurl string) (string, error) {
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
func handleFrontPage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "use http://localhost:5000/api/subscribe?url=<item_url>&email=<mail> to subscribe")
}
func handleSubscribeRequest(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	rawurl := params.Get("url")
	mail := params.Get("email")
	if !CheckEmailSyntax(mail) {
		w.WriteHeader(400)
		fmt.Fprint(w, "Please use correct email")
		return
	}
	id, err := GetItemID(rawurl)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Please use correct direct link")
		return
	}
	verified, err := dbmanager.IsVerified(mail)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "internal database error")
		return
	}
	if verified {
		err := addSubscription(mail, rawurl, id)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "internal error")
			return
		}
		fmt.Fprint(w, "Subscription successful")
		return
	}
	err = startVerification(mail)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "internal error")
		return
	}
	fmt.Fprint(w, "Check your email for verification link")
}
func startVerification(mail string) error {
	key, err := createVerificationKey()
	if err != nil {
		return err
	}
	err = dbmanager.AddVerificationLink(key, mail)
	link := "http://127.0.0.1:5000/verify?key=" + key
	if err != nil {
		return err
	}
	err = netmanager.SendVerificationLink(mail, link)
	return err
}
func createVerificationKey() (string, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(2000000000))
	if err != nil {
		return "", err
	}
	res := fmt.Sprint(num)
	return res, nil
}
func addSubscription(mail, url, id string) error {
	inputChan := make(chan dbmanager.Item, 1)
	outputChan := make(chan netmanager.ParseResult, 1)
	inputChan <- dbmanager.Item{id, -1, url, 1}
	close(inputChan)
	netmanager.GetPrice(inputChan, outputChan)
	res := <-outputChan
	if res.Err != nil {
		return res.Err
	}
	item := dbmanager.Item{id, res.NewPrice, url, 1}
	err := dbmanager.UpdateOrInsertItem(item)
	if err != nil {
		return err
	}
	err = dbmanager.AddSubscription(dbmanager.Subscription{mail, id, true})
	if err != nil {
		return err
	}
	return nil
}
func handleVerification(w http.ResponseWriter, req *http.Request) {
	link := req.URL.Query()
	key := link.Get("key")
	verified, err := dbmanager.Verify(key)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "internal database error")
		return
	}
	if verified {
		fmt.Fprint(w, "Your mail has been verified. You should make a new request for subscription")
		return
	}
	fmt.Fprint(w, "incorrect link")
}
