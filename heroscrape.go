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

var NotComplete = errors.New("Not complete")
var Logger = log.New(ioutil.Discard, "hero-scrape", log.LstdFlags)

type SearchResult struct {
	Image       string
	Title       string
	Description string
}

func (sr *SearchResult) Complete() bool {
	return sr.Title != "" &&
		sr.Image != "" &&
		sr.Description != ""
}

type Strategy interface {
	Scrape(srcUrl *url.URL, doc *goquery.Document) (*SearchResult, error)
}

func Scrape(srcUrl *url.URL, html io.Reader) (*SearchResult, error) {
	return ScrapeWithStrategy(srcUrl, html, NewOgStrategy(), NewHeuristicStrategy())
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

func ScrapeWithStrategy(srcUrl *url.URL, html io.Reader, strategies ...Strategy) (*SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse document")
	}

	var result = new(SearchResult)
	for _, stategy := range strategies {
		finding, err := stategy.Scrape(srcUrl, doc)
		if err != nil {
			return nil, err
		} else if finding != nil {
			Logger.Printf("finding %s \n", srcUrl.String())
			merge(result, finding)
		}

		if result.Complete() {
			Logger.Printf("complete %s \n", srcUrl.String())
			return result, nil
		}
	}

	return result, NotComplete
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

func Debug() {
	Logger = log.New(os.Stderr, "hero-scrape", log.LstdFlags)
}
