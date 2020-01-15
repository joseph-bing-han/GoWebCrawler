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
	//index := time.Now().Nanosecond()
	rand.Seed(time.Now().UnixNano())
	for {
		msg := <-ch
		//fmt.Printf("Thread:%d  Run %s\n", index, num)
		url := msg.url
		id := msg.id
		value, error := cache.Get(url)

		// 判断是否已经查过了
		if error != nil && value.(string) == "" {
			cache.Set(url, 1)
			var className string
			if strings.Contains(url, "www.thewarehouse.co.nz") {
				className = spider.SPIDER_WAREHOUSE
			} else if strings.Contains(url, "www.paknsave.co.nz") {
				className = spider.SPIDER_WAREHOUSE
			}
			spider := spider.Create(className)
			spider.SetURL(url)
			if err := spider.Run(); err == nil {
				mq.Ack(id)
			}
		} else {
			mq.Ack(id)
		}

		time.Sleep(time.Second * time.Duration(2+rand.Intn(8)))
	}
}

func main() {
	cpuNum := runtime.NumCPU() * 2
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
