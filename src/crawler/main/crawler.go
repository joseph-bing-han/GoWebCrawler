package main

import (
	"GoWebCrawler/src/utils/mq"
	"log"
	"time"
)

func main() {

	// 添加2次网址，防止有遗漏的情况

	mq.Add(map[string]interface{}{"url": "https://shop.countdown.co.nz/", "update": false})
	log.Println("Add URL: https://shop.countdown.co.nz/")

	mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/", "update": false})
	log.Println("Add URL: https://www.thewarehouse.co.nz/")

	mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/", "update": false})
	log.Println("Add URL: https://www.kmart.co.nz/")

	//mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/", "update": false})
	//log.Println("Add URL: https://www.paknsaveonline.co.nz/")

	//mq.Add(map[string]interface{}{"url": "https://www.ishopnewworld.co.nz/", "update": false})
	//log.Println("Add URL: https://www.ishopnewworld.co.nz/")

	time.Sleep(time.Second * 20)
	////////////////////////////////////////////////////

	mq.Add(map[string]interface{}{"url": "https://shop.countdown.co.nz/", "update": false})
	log.Println("Add URL: https://shop.countdown.co.nz/")

	mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/", "update": false})
	log.Println("Add URL: https://www.thewarehouse.co.nz/")

	mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/", "update": false})
	log.Println("Add URL: https://www.kmart.co.nz/")

	//mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/", "update": false})
	//log.Println("Add URL: https://www.paknsaveonline.co.nz/")

	//mq.Add(map[string]interface{}{"url": "https://www.ishopnewworld.co.nz/", "update": false})
	//log.Println("Add URL: https://www.ishopnewworld.co.nz/")
}
