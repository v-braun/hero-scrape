package heroscrape_test

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	heroscrape "github.com/v-braun/hero-scrape"
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
	sut := heroscrape.NewHeuristicStrategy()

	u, _ := url.Parse("https://www.google.de")
	res, err := sut.Scrape(u, doc)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, img, res.Image)
}

func scrapeUrl(t *testing.T, pageUrl string) *heroscrape.SearchResult {
	res, err := http.Get(pageUrl)
	if !assert.NoErrorf(t, err, "could not download %s", pageUrl) {
		return nil
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		assert.FailNow(t, "invalid statuscode %d for url %s", res.StatusCode, pageUrl)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if !assert.NoErrorf(t, err, "could not parse %s", pageUrl) {
		return nil
	}

	sut := heroscrape.NewHeuristicStrategy()
	u, _ := url.Parse(pageUrl)

	scrapeRes, err := sut.Scrape(u, doc)
	if !assert.NoErrorf(t, err, "failed scrap %s %v", pageUrl, err) {
		return nil
	}

	return scrapeRes
}

func findByUrl(t *testing.T, pageUrl string) string {
	scrapeRes := scrapeUrl(t, pageUrl)
	if scrapeRes == nil {
		return ""
	}

	return scrapeRes.Image
}

func test(t *testing.T, wg *sync.WaitGroup, pageUrl string, exptected string) {
	wg.Add(1)
	go (func() {
		res := findByUrl(t, pageUrl)
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

	test(t, wg, // relative img url
		"https://rachelbythebay.com/w/2012/10/25/lunch/",
		"https://rachelbythebay.com/w/2012/10/25/lunch/floor.jpg")

	test(t, wg, // abs img url
		"https://www.paulirish.com/2015/advanced-performance-audits-with-devtools/",
		"https://www.paulirish.com/assets/wikipedia-flamechart.jpg")

	test(t, wg, // prio 2 image
		"https://grossmutters-sparstrumpf.de/warum-rohstoffe-nicht-ins-depot-gehoeren/",
		"https://grossmutters-sparstrumpf.de/wp-content/uploads/2018/11/0.jpg")

	test(t, wg, // no img
		"https://akrabat.com/replacing-a-built-in-php-function-when-testing-a-component/",
		"")

	test(t, wg,
		"https://aerotwist.com/blog/cors-for-concern/",
		"https://aerotwist.com/static/blog/cors-for-concern/203-podcast_framed_jpg.jpg")

	test(t, wg,
		"https://nickcraver.com/blog/2018/11/29/stack-overflow-how-we-do-monitoring/",
		"https://nickcraver.com/blog/content/SO-Monitoring/SO-Monitoring-Monitored.png")

	wg.Wait()
}

func TestNegativeCases(t *testing.T) {
	htmlTemplate := `
	<html prefix="og: http://ogp.me/ns#">
	<head>
	</head>
	<body>
	<img src="{{url}}" />
	</body>
	</html>
	`
	heroscrape.Debug()

	// invalid ratio
	u, _ := url.Parse("https://via.placeholder.com/150x1")
	html := strings.Replace(htmlTemplate, "{{url}}", u.String(), 1)
	res, _ := heroscrape.ScrapeWithStrategy(u, strings.NewReader(html), heroscrape.NewHeuristicStrategy())
	assert.Equal(t, "", res.Image)

	// invalid ratio
	u, _ = url.Parse("https://via.placeholder.com/1x150")
	html = strings.Replace(htmlTemplate, "{{url}}", u.String(), 1)
	res, _ = heroscrape.ScrapeWithStrategy(u, strings.NewReader(html), heroscrape.NewHeuristicStrategy())
	assert.Equal(t, "", res.Image)

	// too small
	u, _ = url.Parse("https://via.placeholder.com/10x10")
	html = strings.Replace(htmlTemplate, "{{url}}", u.String(), 1)
	res, _ = heroscrape.ScrapeWithStrategy(u, strings.NewReader(html), heroscrape.NewHeuristicStrategy())
	assert.Equal(t, "", res.Image)

	// bmp
	u, _ = url.Parse("https://www.fileformat.info/format/bmp/sample/43ab63cb34cc4486b09f559a225ce28e/BLK.BMP")
	html = strings.Replace(htmlTemplate, "{{url}}", u.String(), 1)
	res, _ = heroscrape.ScrapeWithStrategy(u, strings.NewReader(html), heroscrape.NewHeuristicStrategy())
	assert.Equal(t, "", res.Image)

	// TIFF file
	u, _ = url.Parse("https://www.fileformat.info/format/tiff/sample/3794038f08df403bb446a97f897c578d/CCITT_1.TIF")
	html = strings.Replace(htmlTemplate, "{{url}}", u.String(), 1)
	res, _ = heroscrape.ScrapeWithStrategy(u, strings.NewReader(html), heroscrape.NewHeuristicStrategy())
	assert.Equal(t, "", res.Image)

	// empty src tag
	u, _ = url.Parse("https://www.fileformat.info/format/bmp/sample/43ab63cb34cc4486b09f559a225ce28e/BLK.BMP")
	html = strings.Replace(htmlTemplate, "{{url}}", "", 1)
	res, _ = heroscrape.ScrapeWithStrategy(u, strings.NewReader(html), heroscrape.NewHeuristicStrategy())
	assert.Equal(t, "", res.Image)

	// invalid src tag
	u, _ = url.Parse("foo://abv/?&$<§819")
	html = strings.Replace(htmlTemplate, "{{url}}", ":?&$<§819", 1)
	res, _ = heroscrape.ScrapeWithStrategy(u, strings.NewReader(html), heroscrape.NewHeuristicStrategy())
	assert.Equal(t, "", res.Image)

}

func TestTitleScrape(t *testing.T) {
	res := scrapeUrl(t, "https://hero-scrape.viktor-braun.de/")
	if assert.NotNil(t, res, "could not find a scrape result") {
		return
	}

	assert.Equal(t, "Google2", res.Title)
}
