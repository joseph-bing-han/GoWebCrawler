package main

import "GoWebCrawler/src/utils/mq"

func main() {
	mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://shop.countdown.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.ishopnewworld.co.nz/"})
}
