package main

import (
	"GoWebCrawler/src/spider"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/mq"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

type Msg struct {
	id  string
	url string
}

func handler(ch chan Msg) {
	rand.Seed(time.Now().UnixNano())
	for {
		time.Sleep(time.Second * time.Duration(12+rand.Intn(8)))
		msg := <-ch
		url := msg.url
		id := msg.id
		//fmt.Println("Read URL:" + url)
		// 判断是否已经查过了
		value, error := cache.Get(url)
		//// todo:test
		//value = ""

		if error != nil && error.Error() == "redis: nil" && value.(string) == "" {
			cache.Set(url, 1)
			var className string
			if strings.Contains(url, "www.thewarehouse.co.nz") {
				className = spider.SPIDER_WAREHOUSE
			} else if strings.Contains(url, "www.paknsaveonline.co.nz") {
				className = spider.SPIDER_PAKNSAVE
			}
			spider := spider.Create(className)
			//log.Println("Get URL: " + url)
			spider.SetURL(url)
			if err := spider.Run(); err == nil {
				mq.Ack(id)
			}
		} else {
			mq.Ack(id)
		}
	}
}

func main() {
	cpuNum := runtime.NumCPU() * 2

	//// todo: test
	//cpuNum = 1

	runtime.GOMAXPROCS(cpuNum)
	ch := make(chan Msg)
	for i := 0; i < cpuNum; i++ {
		go handler(ch)
	}

	// 无限循环，查询消息队列，
	for {
		if id, messages, error := mq.Read(); error == nil {
			ch <- Msg{
				id:  id,
				url: messages["url"].(string),
			}
		}
	}

}
