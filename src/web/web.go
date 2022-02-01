package web

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

// fetches
func Fetch(urlString string) (float64, error) {
	// https://stackoverflow.com/questions/31480710/validate-url-with-standard-package-in-go
	// parse the url
	val, err := url.Parse(urlString)
	if err != nil {
		//TODO do something
		log.Err(err)
		return 0, errors.New("Failed to parse the url string, failed with error: " + err.Error())
	}
	// get the hostname from the url. eg. www.host.com
	host := val.Hostname()

	if strings.Contains(host, "amazon") {
		// run amazon scraper function
		return scrapeAmazon(urlString), nil
	} else if strings.Contains(host, "ebay") {
		// TODO ebay scraper function
		return scrapeEbay(urlString), nil
	}

	return 0, errors.New("url is not currently supported")
}

// scrape an amazon url
func scrapeAmazon(url string) float64 {
	var price float64

	collector := colly.NewCollector()

	collector.OnHTML("span.a-price-whole", func(e *colly.HTMLElement) {
		priceUnClean := e.Text
		// remove commas from price if any
		priceClean := strings.ReplaceAll(priceUnClean, ",", "")
		price, _ = strconv.ParseFloat(priceClean, 64)
	})

	collector.Visit(url)
	collector.Wait()
	return price
}

// scrape an ebay url
func scrapeEbay(url string) float64 {
	var price float64
	collector := colly.NewCollector()

	collector.OnHTML(".mainPrice", func(e *colly.HTMLElement) {
		text := e.ChildAttr("span", "content")
		fmt.Println(text)
		price, _ = strconv.ParseFloat(text, 64)
	})

	collector.Visit(url)
	collector.Wait()

	return price
}

// https://stackoverflow.com/a/48798875
// func trimLeftChar(s string) string {
// 	for i := range s {
// 		if i > 0 {
// 			return s[i:]
// 		}
// 	}
// 	return s[:0]
// }
