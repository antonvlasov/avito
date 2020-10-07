package netmanager

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var key string

func init() {
	var input []byte
	input, err := ioutil.ReadFile("./resources/key.txt")
	if err != nil {
		log.Fatal("could not read key")
	}
	key = (string)(input)
}

func GetPrice(id string) (float64, error) {
	rawurl := "https://www.avito.ru:443/api/15/items/" + id + "?key=" + key
	req, _ := http.NewRequest("GET", rawurl, nil)
	req.Host = "www.avito.ru"
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36 OPR/71.0.3770.198")
	req.Header.Set("accept", `application/json`)

	// Using mitmproxy to avoid getting blocked
	proxyString, _ := url.Parse("http://" + "127.0.0.1:8080")
	PTransport := &http.Transport{
		Proxy:             http.ProxyURL(proxyString),
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		ForceAttemptHTTP2: true,
	}
	client := http.Client{
		Transport: PTransport,
		Timeout:   15 * time.Second,
	}
	//client := http.Client{}
	response, err := client.Do(req)
	data, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return -1, err
	}
	if response.StatusCode != 200 {
		return -1, fmt.Errorf("%v", response.StatusCode)
	}
	return RetrievePrice(data)
}
func RetrievePrice(jsonStream []byte) (float64, error) {
	resp := make(map[string]json.RawMessage)
	err := json.Unmarshal(jsonStream, &resp)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	price := make(map[string]json.RawMessage)
	err = json.Unmarshal(resp["price"], &price)
	var value json.Number
	json.Unmarshal(price["value"], &value)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return value.Float64()
}
