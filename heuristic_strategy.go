package heroscrape

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rubenfonseca/fastimage"
	funk "github.com/thoas/go-funk"
)

var _ Strategy = (*heuristicStrategy)(nil)

var selectorsIgnore = []string{
	"header img",
	"#header img",
	"footer img",
	"#footer img",
	".footer img",
	"#sidebar img",
	".sidebar img",
	"#comment img",
	"#comments img",
}
var selectorsPrio1 = []string{
	"content img",
	".content img",
	"main img",
	".main img",
	"#main img",
	"article img",
	".page img",
	"[role='main'] img",
}

var supportedTypes = []fastimage.ImageType{
	fastimage.JPEG,
	fastimage.PNG,
	fastimage.GIF,
}
var minRatio = 0.6
var maxRatio = 4.0
var minSize = 20000.0

type heuristicStrategy struct {
}

func (og *heuristicStrategy) Scraps(doc *goquery.Document) (*SearchResult, error) {
	allEl := doc.Find("img")
	allEl = allEl.Not(strings.Join(selectorsIgnore, ", "))

	p1El := allEl.Filter(strings.Join(selectorsPrio1, ", "))
	p1Urls := og.getUrls(p1El)
	p1Match := og.findMatch(p1Urls)

	if p1Match != nil {
		return &SearchResult{Image: p1Match.String()}, nil
	}

	// p1Res := checkImages()

	return nil, nil
}
func (og *heuristicStrategy) findMatch(urls []*url.URL) *url.URL {
	for _, u := range urls {
		imgType, size, err := fastimage.DetectImageType(u.String())
		if err != nil {
			continue
		}
		if !og.typeMatch(imgType) {
			continue
		}
		if !og.sizeMatch(size) {
			continue
		}

		return u
	}

	return nil
}

func (og *heuristicStrategy) sizeMatch(s *fastimage.ImageSize) bool {
	w := float64(s.Width)
	h := float64(s.Height)
	ratio := w / h
	size := w * h
	if ratio > maxRatio {
		return false
	} else if ratio < minRatio {
		return false
	} else if size < minSize {
		return false
	} else {
		return true
	}
}

func (og *heuristicStrategy) typeMatch(t fastimage.ImageType) bool {
	return funk.Contains(supportedTypes, t)
}

func (og *heuristicStrategy) getUrls(selections *goquery.Selection) []*url.URL {
	var res []*url.URL
	selections.Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("src", "")
		if href == "" {
			return
		}

		parsedHref, _ := url.Parse(href)
		if parsedHref != nil && parsedHref.IsAbs() {
			res = append(res, parsedHref)
		}
	})

	return res
}
