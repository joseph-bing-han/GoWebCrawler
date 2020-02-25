package main

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/mq"
	"log"
)

func main() {
	count := int32(0)
	model.DB.Table("items").Count(&count)
	log.Println("[INFO]", "Items count =", count)
	for i := 0; i <= int(count/500); i++ {
		var urls []string
		model.DB.Limit(500).Offset(i*500).Table("items").Pluck("url", &urls)
		for _, url := range urls {
			mq.Add(map[string]interface{}{"url": url, "update": true})
			log.Println("[INFO]", "Add URL:", url)
		}
	}

}
