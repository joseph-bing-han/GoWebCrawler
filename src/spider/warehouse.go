package spider

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/conf"
	"GoWebCrawler/src/utils/mq"
	"github.com/bregydoc/gtranslate"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

type Warehouse struct {
	cr       *colly.Collector
	url      string
	isUpdate bool
	tries    int
}

func (w *Warehouse) SetURL(url string, isUpdate bool) {
	if w.cr == nil {
		w.cr = NewCollector(true)
	}
	w.url = url
	w.isUpdate = isUpdate
}

func (w *Warehouse) Run() error {

	if len(w.url) > 0 {

		// 处理所有链接
		w.cr.OnHTML("a[href]", func(e *colly.HTMLElement) {
			if w.isUpdate {
				return
			}
			url := e.Attr("href")
			if strings.Contains(url, "https://www.thewarehouse.co.nz") {
				//fmt.Println(e.Attr("href"))
				checkKey := SPIDER_WAREHOUSE + url
				// todo: test
				if !cache.Has(checkKey) {
					cache.Set(checkKey, 1)
					log.Println("[INFO]", "["+SPIDER_WAREHOUSE+"]", "Get URL: "+url)
					mq.Add(map[string]interface{}{"url": url, "update": false})
				}
			}
		})

		// 处理商品页面数据
		w.cr.OnHTML("div.pdp-main", func(e *colly.HTMLElement) {
			title := e.ChildText(".product-name.hidden-phone")
			if title == "" {
				return
			}

			titleZh := title
			var err error
			titleZh, err = gtranslate.TranslateWithParams(
				title,
				gtranslate.TranslationParams{
					From:  "en",
					To:    "zh",
					Delay: time.Second * 2,
				},
			)
			if err != nil {
				titleZh = title
			}

			itemId := e.ChildAttr("#product-content", "data-itemid")
			price := e.ChildAttr(".pv-price", "data-price")
			productId := e.ChildText("div.row-product-details > div.product-description > div.product-number > span.product-id")
			image := e.ChildAttr(".primary-image", "src")

			url := e.Request.URL.String()

			category, err := cache.Get("Category-" + url)
			if err != nil {
				category = ""
			}

			if productId != "" && price != "" {

				// 在缓存系统中校验是否已经保存过了当天的数据
				checkKey := SPIDER_WAREHOUSE + productId
				if !cache.Has(checkKey) {

					cache.Set(checkKey, 1)
					var item model.Item
					if model.DB.Where("website = ? AND product_id = ?", SPIDER_WAREHOUSE, productId).First(&item).RecordNotFound() {
						// 没找到旧数据时，新建商品记录
						item.Image = image
						item.InternalID = itemId
						item.ProductID = productId
						item.Title = title
						item.TitleZh = titleZh
						item.Website = SPIDER_WAREHOUSE
						item.Url = url
						item.Category = category.(string)
						model.DB.Create(&item)
					} else {
						item.Image = image
						item.InternalID = itemId
						item.ProductID = productId
						item.Title = title
						item.TitleZh = titleZh
						item.Website = SPIDER_WAREHOUSE
						item.Url = url
						item.Category = category.(string)
						model.DB.Save(&item)
					}

					flPrice, _ := strconv.ParseFloat(price, 10)
					model.DB.Model(&item).Association("Prices").Append(model.Price{Price: flPrice})
				}
			}
		})

		w.cr.OnError(func(response *colly.Response, err error) {
			log.Println("[ERROR]", "["+SPIDER_WAREHOUSE+"]", err)
			w.tries--
			if w.tries >= 0 {
				time.Sleep(time.Second)
				response.Request.Retry()
			}

		})

		log.Println("[INFO]", "["+SPIDER_WAREHOUSE+"]", "RUN: "+w.url)
		w.cr.Visit(w.url)

	}
	return nil
}

func init() {
	// 在启动时注册Warehouse类工厂
	Register(SPIDER_WAREHOUSE, func() Spider {
		warehouse := new(Warehouse)
		tries, err := strconv.Atoi(conf.Get("TRIES", "3"))
		if err != nil {
			tries = 3
		}
		warehouse.tries = tries

		return warehouse
	})
}
