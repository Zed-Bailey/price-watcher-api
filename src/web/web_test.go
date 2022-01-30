package web_test

import (
	"RestService/src/web"
	"testing"
)

func TestFetchingPrice(t *testing.T) {
	// NOTE: make sure that the price is upto date with the url
	testUrl := ""
	wantedPrice := 0.0

	price, err := web.Fetch(testUrl)

	if err != nil {
		t.Error(err)
	}

	if price != wantedPrice {
		t.Errorf("Got price %f but wanted %f\n", price, wantedPrice)
	}

}
