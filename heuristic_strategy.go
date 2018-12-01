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
var minRatio = 0.5
var maxRatio = 4.0
var minSize = 20000.0

type heuristicStrategy struct {
}

func NewHeuristicStrategy() Strategy {
	return new(heuristicStrategy)
}

func (hs *heuristicStrategy) Scrape(srcUrl *url.URL, doc *goquery.Document) (*SearchResult, error) {
	allEl := doc.Find("img")
	allEl = allEl.Not(strings.Join(selectorsIgnore, ", "))

	p1El := allEl.Filter(strings.Join(selectorsPrio1, ", "))
	p1Urls := hs.getUrls(p1El, srcUrl)
	p1Match := hs.findMatch(p1Urls)

	if p1Match != nil {
		return &SearchResult{Image: p1Match.String()}, nil
	}

	p2El := allEl
	p2Urls := hs.getUrls(p2El, srcUrl)
	p2Match := hs.findMatch(p2Urls)
	if p2Match != nil {
		return &SearchResult{Image: p2Match.String()}, nil
	}

	// p1Res := checkImages()

	return nil, nil
}
func (hs *heuristicStrategy) findMatch(urls []*url.URL) *url.URL {
	for _, u := range urls {
		imgType, size, err := fastimage.DetectImageType(u.String())

		if err != nil {
			Logger.Printf("fastimage err %v \n", err)
			continue
		}

		Logger.Printf("check type %s \n", u.String())
		if !hs.typeMatch(imgType) {
			Logger.Printf("failed type check %s \n", u.String())
			continue
		}

		Logger.Printf("check size %s \n", u.String())
		if !hs.sizeMatch(size) {
			Logger.Printf("failed size check %s \n", u.String())
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
		Logger.Printf("ratio missmatch | %f > %f  \n", ratio, maxRatio)
		return false
	} else if ratio < minRatio {
		Logger.Printf("ratio missmatch | %f < %f  \n", ratio, minRatio)
		return false
	} else if size < minSize {
		Logger.Printf("size missmatch | %f < %f  \n", size, minSize)
		return false
	} else {
		return true
	}
}

func (og *heuristicStrategy) typeMatch(t fastimage.ImageType) bool {
	return funk.Contains(supportedTypes, t)
}

func (og *heuristicStrategy) getUrls(selections *goquery.Selection, pageUrl *url.URL) []*url.URL {
	var res []*url.URL
	selections.Each(func(i int, s *goquery.Selection) {
		src := s.AttrOr("src", "")
		if src == "" {
			Logger.Printf("no src attr url: %s \n", pageUrl.String())
			return
		}

		parsedHref, err := url.Parse(src)
		if err != nil || parsedHref == nil {
			Logger.Printf("invalid url: %s | src: %s | err: %v | parsed: %v \n", pageUrl.String(), src, err, parsedHref)
			return
		}

		parsedHref = pageUrl.ResolveReference(parsedHref)
		if parsedHref.IsAbs() {
			Logger.Printf("fetched src url: %s | src: %s \n", pageUrl.String(), src)
			res = append(res, parsedHref)
		}
	})

	return res
}
