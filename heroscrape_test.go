package heroscrape_test

import (
	"fmt"
	"net/http"
	"net/url"

	heroscrape "github.com/v-braun/hero-scrape"
)

func ExampleScrape() {
	pageUrl, _ := url.Parse("https://github.com/v-braun/hero-scrape")
	res, _ := http.Get(pageUrl.String())
	defer res.Body.Close()

	result, _ := heroscrape.Scrape(pageUrl, res.Body)
	fmt.Println(result.Image)
}
