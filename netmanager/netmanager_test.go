package netmanager

import (
	"fmt"
	"io/ioutil"
	"testing"
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
