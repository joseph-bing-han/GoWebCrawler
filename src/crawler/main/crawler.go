package main

import "GoWebCrawler/src/utils/mq"

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
	//mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/"})
	//mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/"})
}
