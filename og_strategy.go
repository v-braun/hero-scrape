package heroscrape

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var _ Strategy = (*ogStrategy)(nil)

type ogStrategy struct {
}

func NewOgStrategy() Strategy {
	return new(ogStrategy)
}

func (og *ogStrategy) Scrape(srcUrl *url.URL, doc *goquery.Document) (*SearchResult, error) {
	var result = new(SearchResult)

	result.Title = GetAttrFromSelector(doc, "meta[property='og:title']", "content")
	result.Image = GetAttrFromSelector(doc, "meta[property='og:image']", "content")
	result.Description = GetAttrFromSelector(doc, "meta[property='og:description']", "content")

	return result, nil
}
