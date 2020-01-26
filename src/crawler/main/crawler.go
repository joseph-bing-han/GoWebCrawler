package main

import (
	"GoWebCrawler/src/utils/conf"
	"GoWebCrawler/src/utils/mq"
	"log"
	"net"
	"strings"
	"time"
)

func main() {

	// 向所有Tor网络发命令更换代理ip
	ctlIPs := strings.Split(conf.Get("TOR_PROXY_CTL", ""), ",")
	if len(ctlIPs) > 1 {
		for _, ctlIP := range ctlIPs {
			conn, err := net.Dial("tcp", ctlIP)
			if err != nil {
				log.Println("dial error:", err)
				continue
			}
			time.Sleep(time.Second * 2)
			conn.Write([]byte("AUTHENTICATE \"password\"\r\nSIGNAL NEWNYM\r\n"))
			conn.Close()
		}
	}

	// 添加3次网址，防止有遗漏的情况
	mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://shop.countdown.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.ishopnewworld.co.nz/"})
	time.Sleep(time.Second * 20)
	mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://shop.countdown.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.ishopnewworld.co.nz/"})
	time.Sleep(time.Second * 20)
	mq.Add(map[string]interface{}{"url": "https://www.thewarehouse.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.paknsaveonline.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.kmart.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://shop.countdown.co.nz/"})
	mq.Add(map[string]interface{}{"url": "https://www.ishopnewworld.co.nz/"})
}
