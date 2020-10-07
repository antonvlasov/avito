package netmanager

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
)

var usr, pswd string

func init() {
	var input []byte
	input, err := ioutil.ReadFile("./resources/mail_data.txt")
	if err != nil {
		log.Fatal("could not read mail data")
	}
	s := strings.Split((string)(input), " ")
	if len(s) != 2 {
		log.Fatal("could not read mail data")
	}
	usr = s[0]
	pswd = s[1]
}
func Notify(mail, url string, price float64) {
	ps := fmt.Sprintf("%v", price)
	auth := smtp.PlainAuth("", usr, pswd, "smtp.gmail.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{mail}
	msg := []byte("To: " + mail + "\r\n" +
		"Subject: Change of price\r\n" +
		"\r\n" +
		"Price on item: " + url + " now costs " + ps + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:25", auth, usr, to, msg)
	if err != nil {
		fmt.Println(err)
	}
}
