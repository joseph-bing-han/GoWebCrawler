package main

import (
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/conf"
	"GoWebCrawler/src/utils/mq"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"net"
	"strings"
	"time"
)

func newCollector() *colly.Collector {
	cr := colly.NewCollector()

	// 使用随机User Agent
	extensions.RandomUserAgent(cr)

	// 添加referer防盗链
	extensions.Referer(cr)

	cr.SetRequestTimeout(time.Minute * 2)
	return cr
}

func main() {

	//// 向所有Tor网络发命令更换代理ip
	ctlIPs := strings.Split(conf.Get("TOR_PROXY_CTL", ""), ",")
	proxyIPs := strings.Split(conf.Get("TOR_PROXY", ""), ",")
	if len(ctlIPs) > 1 {
		ch := make(chan int, len(ctlIPs))
		for index, ctlIP := range ctlIPs {
			conn, err := net.Dial("tcp", ctlIP)
			if err != nil {
				log.Println("dial error:", err)
				continue
			}
			time.Sleep(time.Second * 2)
			conn.Write([]byte("AUTHENTICATE \"password\"\r\nSIGNAL NEWNYM\r\n"))
			conn.Close()
			log.Println(ctlIP, "Refresh Proxy Server IP.")
			go func(proxy string, ctrl string) {
				cr := newCollector()
				cr.UserAgent = "Mozilla/5.0 (Windows NT 6.1; rv:60.0) Gecko/20100101 Firefox/60.0"
				cr.SetProxy(proxy)
				cr.OnResponse(func(response *colly.Response) {
					ch <- 1
				})
				cr.OnResponse(func(response *colly.Response) {
					cookie := response.Headers.Get("Set-Cookie")
					fmt.Println("SAVE COOKIE", "paknsave-cookie-key-"+proxy, cookie)
					cache.Set("paknsave-cookie-key-"+proxy, cookie)

				})
				cr.OnError(func(response *colly.Response, err error) {
					log.Println("ERROR", ctrl, err)
					conn, err1 := net.Dial("tcp", ctrl)
					if err1 != nil {
						log.Println("ERROR", ctrl, err1)
					}
					time.Sleep(time.Second * 2)
					conn.Write([]byte("AUTHENTICATE \"password\"\r\nSIGNAL NEWNYM\r\n"))
					conn.Close()
					log.Println(ctrl, "Refresh Proxy Server IP.")
					time.Sleep(time.Second * 5)
					response.Request.Retry()
				})
				cr.Visit("https://www.paknsaveonline.co.nz/")
			}(proxyIPs[index], ctlIP)

		}
		<-ch
	}

	// 添加3次网址，防止有遗漏的情况
	//mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/"})
	//log.Println("Add URL: https://www.thewarehouse.co.nz/")
	//
	mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/"})
	log.Println("Add URL: https://www.paknsaveonline.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/"})
	//log.Println("Add URL: https://www.kmart.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://shop.countdown.co.nz/"})
	//log.Println("Add URL: https://shop.countdown.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://www.ishopnewworld.co.nz/"})
	//log.Println("Add URL: https://www.ishopnewworld.co.nz/")

	//time.Sleep(time.Second * 20)
	////////////////////////////////////////////////////
	//
	//mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/"})
	//log.Println("Add URL: https://www.thewarehouse.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/"})
	//log.Println("Add URL: https://www.paknsaveonline.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/"})
	//log.Println("Add URL: https://www.kmart.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://shop.countdown.co.nz/"})
	//log.Println("Add URL: https://shop.countdown.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://www.ishopnewworld.co.nz/"})
	//log.Println("Add URL: https://www.ishopnewworld.co.nz/")
	//
	//time.Sleep(time.Second * 20)
	//////////////////////////////////////////////////
	//
	//mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/"})
	//log.Println("Add URL: https://www.thewarehouse.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/"})
	//log.Println("Add URL: https://www.paknsaveonline.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/"})
	//log.Println("Add URL: https://www.kmart.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://shop.countdown.co.nz/"})
	//log.Println("Add URL: https://shop.countdown.co.nz/")
	//
	//mq.Add(map[string]interface{}{"url": "https://www.ishopnewworld.co.nz/"})
	//log.Println("Add URL: https://www.ishopnewworld.co.nz/")

}
