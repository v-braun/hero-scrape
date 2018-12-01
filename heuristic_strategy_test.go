package heroscrape

import (
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestHeuristicStrategyFilterDeadZones(t *testing.T) {
	img := "https://www.google.de/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"
	raw := `
	<html prefix="og: http://ogp.me/ns#">
	<body>
		<img id="img1" />
		<div id="footer"><img id="img2" /></div>
		<article role="main">
		<img id="imgX" src="` + img + `" />
		</article>
		<div class="footer"><div class="inner"><img id="img3" /></div></div>
		<footer><img id="img4" /></footer>
	</body>
	</html>	
	`
	html := strings.NewReader(raw)
	doc, _ := goquery.NewDocumentFromReader(html)
	sut := new(heuristicStrategy)

	res, err := sut.Scraps(doc)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, img, res.Image)
}

func findByUrl(t *testing.T, url string) string {
	res, err := http.Get(url)
	if !assert.NoErrorf(t, err, "could not download %s", url) {
		return ""
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		assert.FailNow(t, "invalid statuscode %d for url %s", res.StatusCode, url)
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if !assert.NoErrorf(t, err, "could not parse %s", url) {
		return ""
	}

	sut := new(heuristicStrategy)
	scrapeRes, err := sut.Scraps(doc)
	if !assert.NoErrorf(t, err, "failed scrap %s %v", url, err) {
		return ""
	}

	if !assert.NotNil(t, scrapeRes) {
		return ""
	}

	return scrapeRes.Image
}

func test(t *testing.T, wg *sync.WaitGroup, url string, exptected string) {
	wg.Add(1)
	go (func() {
		res := findByUrl(t, url)
		assert.Equal(t, exptected, res)
		wg.Done()
	})()

}

func TestHeuristicStrategyBlogs(t *testing.T) {
	wg := &sync.WaitGroup{}
	test(t, wg,
		"https://blog.sindresorhus.com/my-macos-10-14-wishlist-c499448afdd6",
		"https://cdn-images-1.medium.com/max/1600/1*e8JW87HyeIWFlACUmLjDOg.jpeg")

	test(t, wg,
		"https://blog.sindresorhus.com/gifski-972692460aa5",
		"https://cdn-images-1.medium.com/max/2000/1*9g6fkWCL2xylg7moinRWVQ.png")

	test(t, wg,
		"https://blog.ghost.org/2-0/",
		"https://blog.ghost.org/content/images/2018/08/editor.png")

	wg.Wait()
}
