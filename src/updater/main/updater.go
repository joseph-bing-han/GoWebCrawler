package main

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/mq"
	"log"
)

func main() {
	count := int32(0)
	model.DB.Table("sources").Count(&count)
	log.Println("[INFO]", "[Updater]", "Sources count =", count)
	for i := 0; i <= int(count/500); i++ {
		var sources []model.Source
		model.DB.Limit(500).Offset(i*500).Where("active = ?", 1).Find(&sources)
		for _, source := range sources {
			mq.Add(map[string]interface{}{"url": source.Url, "update": true})
			cache.Set("Category-"+source.Url, source.Category)
			log.Println("[INFO]", "[Updater]", "Add URL:", source.Url)
		}
	}
}
