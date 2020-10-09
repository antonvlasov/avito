package netmanager

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
)

var usr, pswd string

func init() {
	path, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	var input []byte
	input, err = ioutil.ReadFile("../resources/mail_data.txt")
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
func Notify(mail, url string, price float64) error {
	var ps string
	if price < 0 {
		ps = "has been deleted"
	} else {
		ps = fmt.Sprintf("now costs %v", price)
	}
	message := "Item: " + url + ps
	return sendMail(mail, message)
}
func SendVerificationLink(mail, link string) error {
	message := fmt.Sprint("To verify your mail for avito price subscription please follow the link: ", link)
	return sendMail(mail, message)
}
func sendMail(reciever string, message string) error {
	auth := smtp.PlainAuth("", usr, pswd, "smtp.gmail.com")
	to := []string{reciever}
	msg := []byte("To: " + reciever + "\r\n" +
		"Subject: Price subscription\r\n" +
		"\r\n" +
		message + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:25", auth, usr, to, msg)
	if err != nil {
		return err
	}
	return nil
}
