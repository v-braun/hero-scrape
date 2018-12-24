package heroscrape

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var _ Strategy = (*ogStrategy)(nil)

type ogStrategy struct {
}

// NewOgStrategy returns a new Strategy that search for OG meta tags
func NewOgStrategy() Strategy {
	return new(ogStrategy)
}

func (og *ogStrategy) Scrape(srcURL *url.URL, doc *goquery.Document) (*SearchResult, error) {
	var result = new(SearchResult)

	result.Title = GetAttrFromSelector(doc, "meta[property='og:title']", "content")
	result.Image = GetAttrFromSelector(doc, "meta[property='og:image']", "content")
	result.Description = GetAttrFromSelector(doc, "meta[property='og:description']", "content")

	return result, nil
}
