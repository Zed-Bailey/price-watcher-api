package web_test

import (
	"RestService/src/web"
	"testing"
)

func TestFetchingPriceAmazon(t *testing.T) {
	// NOTE: make sure that the price is upto date with the url
	// at the time of writing this t est the price was $1349 AUD for this ipad pro
	testUrl := "https://www.amazon.com.au/Apple-11-inch-iPad-Pro-Wi-Fi-256GB/dp/B093RLHVQ4/ref=sr_1_5?crid=1MCFI784VLJ1L&keywords=ipad%2Bpro&qid=1643599137&sprefix=ipad%2Bpro%2Caps%2C291&sr=8-5&th=1"
	wantedPrice := 1349.0

	got, err := web.Fetch(testUrl)
	if err != nil {
		t.Error(err)
	}

	if got != wantedPrice {
		t.Errorf("AMAZON: Got price %f but wanted %f\n", got, wantedPrice)
	}

}

func TestFetchingPriceEbay(t *testing.T) {
	// NOTE: make sure that the price is upto date with the url
	// at the time of writing this test the price is $1489
	testUrl := "https://www.ebay.com.au/itm/254795913956?var=&hash=item3b53051ae4"
	wantedPrice := 1489.0

	got, err := web.Fetch(testUrl)
	if err != nil {
		t.Error(err)
	}

	if got != wantedPrice {
		t.Errorf("EBAY: Got price %f but wanted %f\n", got, wantedPrice)
	}

}
