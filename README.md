# hero-scrape
> Find the hero (main) image of an URL 

[![Build Status](https://travis-ci.org/v-braun/hero-scrape.svg?branch=master)](https://travis-ci.org/v-braun/hero-scrape)
[![codecov](https://codecov.io/gh/v-braun/hero-scrape/branch/master/graph/badge.svg)](https://codecov.io/gh/v-braun/hero-scrape)

By [v-braun - viktor-braun.de](https://viktor-braun.de).

<p align="center">
<img width="70%" src="https://raw.githubusercontent.com/v-braun/hero-scrape/master/logo.png?sanitize=true" />
</p>

## Demo
See a demo on https://hero-scrape.viktor-braun.de

## Description
hero-scrape extracts the main image of a webpage.
It use different strategies to find the main images (OpenGraph HTML Tags and heuristic search).
You can use the existing strategies or implement your own.

To find the "biggest" image it is necessary to download it. [fastimage](https://github.com/rubenfonseca/fastimage/) is the perfect choice for that job.




## Installation

```bash
go get github.com/v-braun/hero-scrape
```

## Usage

**With pre configured strategies**

```go
pageUrl, _ := url.Parse("https://github.com/v-braun/hero-scrape")
res, _ := http.Get(pageUrl.String())
defer res.Body.Close()

result, _ := heroscrape.Scrape(pageUrl, res.Body)
fmt.Println(result.Image)

```

**With cusom strategies**

```go
pageUrl, _ := url.Parse("https://github.com/v-braun/hero-scrape")
res, _ := http.Get(pageUrl.String())
defer res.Body.Close()

result, _ := heroscrape.ScrapeWithStrategy(pageUrl, res.Body, , NewOgStrategy(), NewHeuristicStrategy(), YourOwnStrategy())
fmt.Println(result.Image)

```


## Related Projects
- [hero-scrape] (https://github.com/v-braun/hero-scrape-web) Demo for this lib
- [fastimage](https://github.com/rubenfonseca/fastimage/) Finds the type and/or size of a remote image given its uri, by fetching as little as needed.
- [goquery](https://github.com/PuerkitoBio/goquery) A little like that j-thing, only in Go.

## Known Issues

If you discover any bugs, feel free to create an issue on GitHub fork and
send me a pull request.

[Issues List](https://github.com/v-braun/hero-scrape/issues).

## Authors

![image](https://avatars3.githubusercontent.com/u/4738210?v=3&s=50)  
[v-braun](https://github.com/v-braun/)



## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


## License

See [LICENSE](https://github.com/v-braun/hero-scrape/blob/master/LICENSE).
