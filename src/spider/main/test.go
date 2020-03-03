package main

import "GoWebCrawler/src/spider"

func main() {
	spider := spider.Create(spider.SPIDER_PAKNSAVE)
	//log.Println("Get URL: " + url)
	spider.SetURL("https://www.paknsaveonline.co.nz/product/5045842_kgm_000pns?name=brushed-potatoes", true)
	spider.Run()
}
