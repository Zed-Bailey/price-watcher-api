package web

import (
	"errors"
	"net/url"
	"strings"

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
	} 
	else if strings.Contains(host, "ebay") {
		// TODO ebay scraper function
		return scrapeEbay(urlString), nil
	}
	

	return 0, errors.New("url is not currently supported")
}

// scrape an amazon url
func scrapeAmazon(url string) float64 {
	// TODO implement
	return 0
}

// scrape an ebay url
func scrapeEbay(url string) float64 {
	// TODO implement
	return 0
}