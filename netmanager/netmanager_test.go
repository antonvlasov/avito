package netmanager

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/antonvlasov/avito/dbmanager"
)

func TestRetrievePrice(t *testing.T) {
	var input []byte
	input, err := ioutil.ReadFile("../testdata/response.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	result, _ := RetrievePrice(input)
	res := 800.0
	if result != res {
		t.Errorf("RetrievePrice failed on item, expected price = %v, got %v", res, result)
	}

}
func TestParse(t *testing.T) {
	for i := 0; i < 2; i++ {
		items := []dbmanager.Item{
			{"wrongid", 104, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya._kanonicheskie_pisaniya_wrongid", 1},
			{"1867132320", 400, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya_kanonicheskaya_1867132320", 0},
			{"1897268986", 1500, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya._kanonicheskie_pisaniya_1897268986", 0},
			{"1897268986", 1000, "https://www.avito.ru/moskva/knigi_i_zhurnaly/bibliya._kanonicheskie_pisaniya_1897268986", 1},
			dbmanager.Item{},
		}
		chanelCount := 100
		parseOutputChan := make(chan ParseResult)
		parseInputChan := make(chan dbmanager.Item, len(items))
		for i := range items {
			parseInputChan <- items[i]
		}
		close(parseInputChan)
		for i := 0; i < chanelCount; i++ {
			go GetPrice(parseInputChan, parseOutputChan)
		}
		for range items {
			res := <-parseOutputChan
			fmt.Println(res)
		}
		proxy = ""
	}
}
func TestNotify(t *testing.T) {
	err := Notify("noreplyavitobot@gmail.com", "test_item", 200.0)
	if err != nil {
		t.Errorf("Notify failed with err %v", err)
	}
	err = Notify("mail", "test_item", -1)
	if err == nil {
		t.Errorf("Notify failed, no error in invalid mail")
	}
}
