package main

import "GoWebCrawler/src/spider"

func main() {
	spider := spider.Create(spider.SPIDER_COUNTDOWN)
	//log.Println("Get URL: " + url)
	spider.SetURL("https://shop.countdown.co.nz/", false)
	spider.Run()
}
