package spider

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/mq"
	"github.com/bregydoc/gtranslate"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

type Warehouse struct {
	cr  *colly.Collector
	url string
}

func (w *Warehouse) SetURL(url string) {
	if w.cr == nil {
		w.cr = NewCollector(true)
	}
	w.url = url
}

func (w *Warehouse) Run() error {

	if len(w.url) > 0 {

		// 处理所有链接
		w.cr.OnHTML("a[href]", func(e *colly.HTMLElement) {
			url := e.Attr("href")
			if strings.Contains(url, "https://www.thewarehouse.co.nz") {
				//fmt.Println(e.Attr("href"))

				// todo
				if !cache.Has(url)  {
					cache.Set(url, 1)
					mq.Add(map[string]interface{}{"url": url})
				}
			}
		})

		// 处理商品页面数据
		w.cr.OnHTML("div.pdp-main", func(e *colly.HTMLElement) {
			title := e.ChildText(".product-name.hidden-phone")
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
			itemId := e.ChildAttr("#product-content", "data-itemid")
			price := e.ChildAttr(".pv-price", "data-price")
			productId := e.ChildText("div.row-product-details > div.product-description > div.product-number > span.product-id")
			image := e.ChildAttr(".primary-image", "src")
			if productId != "" && price != "" {

				// 在缓存系统中校验是否已经保存过了当天的数据
				checkKey := time.Now().Format("20060102") + SPIDER_WAREHOUSE + productId
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
						model.DB.Create(&item)
					}

					flPrice, _ := strconv.ParseFloat(price, 10)
					model.DB.Model(&item).Association("Prices").Append(model.Price{Price: flPrice})
				}
			}
		})
		log.Println("Warehouse Run: " + w.url)
		w.cr.Visit(w.url)

	}
	return nil
}

func init() {
	// 在启动时注册Warehouse类工厂
	Register(SPIDER_WAREHOUSE, func() Spider {
		return new(Warehouse)
	})
}
