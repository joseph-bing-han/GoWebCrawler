package main

import (
	"GoWebCrawler/src/utils/mq"
	"fmt"
)

func main() {
	fmt.Println("Clear All Message Queue")
	mq.Clear()
}
