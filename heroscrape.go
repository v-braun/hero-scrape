package heroscrape

import (
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type ImageLocation string

var NotComplete = errors.New("Not complete")

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
	Scraps(doc *goquery.Document) (*SearchResult, error)
}

func Scrap(html io.Reader) (*SearchResult, error) {
	return ScrapWithStrategy(html, NewOgStrategy())
}

func ScrapWithStrategy(html io.Reader, strategies ...Strategy) (*SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse document")
	}

	var result = new(SearchResult)
	for _, stategy := range strategies {
		finding, err := stategy.Scraps(doc)
		if err != nil {
			return nil, err
		} else if finding != nil {
			merge(result, finding)
		}

		if result.Complete() {
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
