package main

import (
	"GoWebCrawler/src/utils/conf"
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
				log.Println("[ERROR]", err)
				continue
			}
			time.Sleep(time.Second * 2)
			conn.Write([]byte("AUTHENTICATE \"password\"\r\nSIGNAL NEWNYM\r\n"))
			conn.Close()
			log.Println(ctlIP, "Refresh Proxy Server IP.")
		}
	}
}
