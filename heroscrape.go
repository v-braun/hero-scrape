package heroscrape

import (
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type ImageLocation string

// ErrNotComplete will be returned if the Scrape was not completley done
var ErrNotComplete = errors.New("Not complete")

// Logger instance for the entire module
var Logger = log.New(ioutil.Discard, "hero-scrape", log.LstdFlags)

// SearchResult represents the scrape result
type SearchResult struct {
	Image       string
	Title       string
	Description string
}

// Complete returns true if the SearchResult has found everything
func (sr *SearchResult) Complete() bool {
	return sr.Title != "" &&
		sr.Image != "" &&
		sr.Description != ""
}

// Strategy interface represents an interface for scraping an website
type Strategy interface {
	Scrape(srcURL *url.URL, doc *goquery.Document) (*SearchResult, error)
}

// Scrape the given url
func Scrape(srcURL *url.URL, html io.Reader) (*SearchResult, error) {
	return ScrapeWithStrategy(srcURL, html, NewOgStrategy(), NewHeuristicStrategy())
}

// TODO
// func ScrapeUrl(srcUrl *url.URL) (*SearchResult, error) {
// 	res, err := http.Get(pageUrl)
// 	if err != nil{
// 		return nil, err
// 	}

// 	defer res.Body.Close()
// 	return ScrapeWithStrategy(srcUrl, html, NewOgStrategy(), NewHeuristicStrategy())
// }

// ScrapeWithStrategy scrapes the given url with the given strategy
func ScrapeWithStrategy(srcURL *url.URL, html io.Reader, strategies ...Strategy) (*SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse document")
	}

	var result = new(SearchResult)
	for _, stategy := range strategies {
		finding, err := stategy.Scrape(srcURL, doc)
		if err != nil {
			return nil, err
		} else if finding != nil {
			Logger.Printf("finding %s \n", srcURL.String())
			merge(result, finding)
		}

		if result.Complete() {
			Logger.Printf("complete %s \n", srcURL.String())
			return result, nil
		}
	}

	return result, ErrNotComplete
}

func merge(dest *SearchResult, src *SearchResult) {
	if dest.Title == "" {
		dest.Title = src.Title
	}
	if dest.Image == "" {
		dest.Image = src.Image
	}
	if dest.Description == "" {
		dest.Description = src.Description
	}
}

// Debug enables the module log debugging
func Debug() {
	Logger = log.New(os.Stderr, "hero-scrape", log.LstdFlags)
}
