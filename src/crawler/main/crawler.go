package main

import (
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/mq"
	"fmt"
)

func main() {
	//c := colly.NewCollector()
	//
	//// Find and visit all links
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	e.Request.Visit(e.Attr("href"))
	//})
	//
	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL)
	//})
	//
	//c.Visit("http://go-colly.org/")
	cache.Set("abc", "25467")
	fmt.Println("finish")
	mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/"})
	fmt.Println(cache.Get("abc"))
}
