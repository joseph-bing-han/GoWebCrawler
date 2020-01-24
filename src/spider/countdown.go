package spider

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/mq"
	"github.com/bregydoc/gtranslate"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Countdown struct {
	cr  *colly.Collector
	url string
}

func (w *Countdown) SetURL(url string) {
	if w.cr == nil {
		w.cr = NewCollector(true)
	}
	w.url = url
}

func (w *Countdown) Run() error {

	if len(w.url) > 0 {

		// 处理所有链接
		w.cr.OnHTML("a[href]", func(e *colly.HTMLElement) {
			url := e.Attr("href")
			//fmt.Println("Get URL:" + url)
			if match, _ := regexp.MatchString(`^/[\w\W]+$`, url); match {
				if strings.Contains(url, "/shop/recipe") {
					return
				}
				url = "https://shop.countdown.co.nz" + url
				checkKey := time.Now().Format("20060102") + SPIDER_COUNTDOWN + url
				// todo: test
				if !cache.Has(checkKey) {
					cache.Set(checkKey, 1)
					//fmt.Println("Add URL: " + url)
					mq.Add(map[string]interface{}{"url": url})
				}
			}
		})

		// 处理商品页面数据
		w.cr.OnHTML("body", func(e *colly.HTMLElement) {
			productId := e.ChildAttr("input[name='stockcode']", "value")

			if productId == "" {
				return
			}

			title := e.ChildText("h1")
			if title == "" {
				return
			}

			titleZh, error := gtranslate.TranslateWithParams(
				title,
				gtranslate.TranslationParams{
					From: "en",
					To:   "zh",
				},
			)
			if error != nil {
				titleZh = title
			}

			data := strings.Split(e.ChildText("span.price"), " ")
			price := ""
			unit := ""
			if len(data) > 1 {
				price = strings.Replace(data[0], "$", "", -1)
				unit = data[1]
			}
			image := "https://shop.countdown.co.nz" + e.ChildAttr("img.product-image", "src")

			//fmt.Println(title + " > " + productId + " > " + price + "/" + unit + " ---> " + image)

			if productId != "" && price != "" {
				// 在缓存系统中校验是否已经保存过了当天的数据
				checkKey := time.Now().Format("20060102") + SPIDER_COUNTDOWN + productId
				if !cache.Has(checkKey) {

					cache.Set(checkKey, 1)
					var item model.Item
					if model.DB.Where("website = ? AND product_id = ?", SPIDER_COUNTDOWN, productId).First(&item).RecordNotFound() {
						// 没找到旧数据时，新建商品记录
						item.Image = image
						item.Unit = unit
						item.ProductID = productId
						item.Title = title
						item.TitleZh = titleZh
						item.Website = SPIDER_COUNTDOWN
						model.DB.Create(&item)
					}

					flPrice, _ := strconv.ParseFloat(price, 10)
					model.DB.Model(&item).Association("Prices").Append(model.Price{Price: flPrice})
				}
			}
		})
		log.Println("Countdown Run: " + w.url)
		w.cr.Visit(w.url)

	}
	return nil
}

func init() {
	// 在启动时注册Countdown类工厂
	Register(SPIDER_COUNTDOWN, func() Spider {
		return new(Countdown)
	})
}
