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
		// todo: test
		time.Sleep(time.Second * time.Duration(12+rand.Intn(8)))

		msg := <-ch
		url := msg.url
		id := msg.id
		//fmt.Println("Read URL:" + url)

		// todo: test
		//if true {
		// 判断是否已经查过了
		checkKey := time.Now().Format("20060102") + "RUN" + url
		if !cache.Has(checkKey) {
			cache.Set(checkKey, 1)
			var className string
			if strings.Contains(url, "www.thewarehouse.co.nz") {
				className = spider.SPIDER_WAREHOUSE
			} else if strings.Contains(url, "www.paknsaveonline.co.nz") {
				className = spider.SPIDER_PAKNSAVE
			} else if strings.Contains(url, "www.kmart.co.nz") {
				className = spider.SPIDER_KMART
			} else if strings.Contains(url, "shop.countdown.co.nz") {
				className = spider.SPIDER_COUNTDOWN
			} else if strings.Contains(url, "www.ishopnewworld.co.nz") {
				className = spider.SPIDER_NEWWORLD
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

	// todo: test
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
