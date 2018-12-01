package heroscrape

import "github.com/PuerkitoBio/goquery"

func GetAttrFromSelector(doc *goquery.Document, selector string, attrName string) string {
	var tag = doc.Find(selector).First()
	return tag.AttrOr(attrName, "")
}
