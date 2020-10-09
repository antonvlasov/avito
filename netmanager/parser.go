package netmanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antonvlasov/avito/dbmanager"
)

var key string
var proxy string = ""

//using proxy to aviod getting blocked
//var proxy string = "http://127.0.0.1:8080"

func init() {
	var input []byte
	input, err := ioutil.ReadFile("../resources/key.txt")
	if err != nil {
		log.Fatal("could not read key")
	}
	key = (string)(input)
}

type ParseResult struct {
	dbmanager.Item
	NewPrice float64
	Err      error
}

func GetPrice(inputChan <-chan dbmanager.Item, outputChan chan<- ParseResult) {
	for item, open := <-inputChan; open; item, open = <-inputChan {
		var err error
		time.Sleep(time.Duration(rand.Intn(500)+1000) * time.Millisecond)
		req, err := createRequest(item.Id)
		if err != nil {
			outputChan <- ParseResult{Err: err}
			continue
		}
		// Using mitmproxy to avoid getting blocked
		client, err := createClient(proxy)
		//client, err := createClient("")
		if err != nil {
			outputChan <- ParseResult{Err: err}
			continue
		}
		response, err := client.Do(req)
		if err != nil {
			outputChan <- ParseResult{Err: err}
			continue
		}
		newPrice, err := readResponse(response)
		//newPrice, err := 0.0, nil
		if err != nil {
			outputChan <- ParseResult{Err: err}
			continue
		}
		outputChan <- ParseResult{Item: item, NewPrice: newPrice, Err: nil}
	}
}
func readResponse(response *http.Response) (newPrice float64, err error) {
	data, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		if response.StatusCode == 400 {
			return -1, nil
		}
		fmt.Printf("%s \n", data)
		return 0, fmt.Errorf("%v", response.StatusCode)
	}
	return RetrievePrice(data)
}
func createRequest(id string) (req *http.Request, err error) {
	rawurl := "https://www.avito.ru:443/api/15/items/" + id + "?key=" + key
	req, err = http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return req, err
	}
	req.Host = "www.avito.ru"
	//req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; America Online Browser 1.1; Windows NT 5.0)")
	return
}
func createClient(proxy string) (client *http.Client, err error) {
	if proxy == "" {
		return http.DefaultClient, nil
	}
	proxyString, err := url.Parse(proxy)
	if err != nil {
		return
	}
	PTransport := &http.Transport{
		Proxy:             http.ProxyURL(proxyString),
		ForceAttemptHTTP2: true,
	}
	client = &http.Client{
		Transport: PTransport,
		Timeout:   15 * time.Second,
	}
	return client, nil
}

func RetrievePrice(jsonStream []byte) (float64, error) {
	resp := make(map[string]json.RawMessage)
	err := json.Unmarshal(jsonStream, &resp)
	if err != nil {
		return 0, err
	}
	price := make(map[string]string)
	err = json.Unmarshal(resp["price"], &price)
	if err != nil {
		return 0, err
	}
	var value string = strings.Replace(price["value"], " ", "", -1)
	return strconv.ParseFloat(value, 64)
}
