package heroscrape_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	heroscrape "github.com/v-braun/hero-scrape"
)

var fullHtml string = `
<html prefix="og: http://ogp.me/ns#">
<head>
<title>The Rock (1996)</title>
<meta property="og:title" content="The Rock" />
<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />
<meta property="og:description" content="A movie"/>
</head>
</html>
`

var nonFullHtml string = `
<html prefix="og: http://ogp.me/ns#">
<head>
<title>The Rock (1996)</title>
<meta property="og:title" content="The Rock" />
<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />
</head>
</html>
`

func TestOgStrategy(t *testing.T) {
	u, _ := url.Parse("http://www.imdb.com")
	html := strings.NewReader(fullHtml)
	result, err := heroscrape.ScrapeWithStrategy(u, html, heroscrape.NewOgStrategy())
	assert.NoError(t, err)

	assert.Equal(t, "The Rock", result.Title)
	assert.Equal(t, "http://ia.media-imdb.com/images/rock.jpg", result.Image)
	assert.Equal(t, "A movie", result.Description)
}

func TestOgStrategyPartial(t *testing.T) {
	u, _ := url.Parse("http://www.imdb.com")
	html := strings.NewReader(nonFullHtml)
	result, err := heroscrape.ScrapeWithStrategy(u, html, heroscrape.NewOgStrategy())
	assert.Equal(t, heroscrape.ErrNotComplete, err)

	assert.Equal(t, "The Rock", result.Title)
	assert.Equal(t, "http://ia.media-imdb.com/images/rock.jpg", result.Image)
	assert.Equal(t, "", result.Description)
}
