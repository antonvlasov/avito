package httpserver

import (
	"testing"
	"time"

	"github.com/antonvlasov/avito/handler"
)

func TestMailCheck(t *testing.T) {
	mails := []string{
		"vlasovanto.n@bk.ru",
		"noreplyavitobot@gmail.com",
		"notamail",
		"",
	}
	answers := []bool{true, true, false, false}
	for i := range mails {
		res := CheckEmailSyntax(mails[i])
		if res != answers[i] {
			t.Errorf("not correct result for mail %v", mails[i])
		}
	}
}
func TestCreateLink(t *testing.T) {
	for i := 0; i < 20; i++ {
		_, err := createVerificationKey()
		if err != nil {
			t.Errorf("createVerificationLink eror %v", err)
		}
	}
}
func TestRun(t *testing.T) {
	handler.Run(10)
	go Run()
	timer := time.NewTimer(120 * time.Second)
	<-timer.C
	handler.Stop()
}
