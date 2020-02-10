package spider

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/mq"
	"encoding/json"
	"github.com/bregydoc/gtranslate"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"time"
)

type Paknsave struct {
	cr  *colly.Collector
	url string
}

func (w *Paknsave) SetURL(url string) {
	if w.cr == nil {
		w.cr = NewCollector(false)
	}
	w.url = url
}

func (w *Paknsave) Run() error {

	if len(w.url) > 0 {

		// 处理所有链接
		w.cr.OnHTML("a[href]", func(e *colly.HTMLElement) {
			url := e.Attr("href")
			//fmt.Println("Get URL:" + url)
			if match, _ := regexp.MatchString(`^/[\w\W]+$`, url); match {
				url = "https://www.paknsaveonline.co.nz" + url
				checkKey := time.Now().Format("20060102") + SPIDER_PAKNSAVE + url
				// todo: test
				//value = ""
				if !cache.Has(checkKey) {
					cache.Set(checkKey, 1)
					//log.Println("Add URL: " + url)
					mq.Add(map[string]interface{}{"url": url})
				}
			}
		})

		// 处理商品页面数据
		w.cr.OnHTML(".fs-product-detail", func(e *colly.HTMLElement) {
			title := e.ChildText("h1")
			if title == "" {
				return
			}

			titleZh, error := gtranslate.TranslateWithParams(
				title,
				gtranslate.TranslationParams{
					From: "en",
					To:   "zh",
					Delay: time.Second * 2,
				},
			)
			if error != nil {
				titleZh = title
			}
			productId := ""
			optionsJson := e.ChildAttr("div.fs-product-detail__wishlist", "data-options")
			if len(optionsJson) > 0 {
				var options map[string]interface{}
				json.Unmarshal([]byte(optionsJson), &options)
				productId = options["productId"].(string)
			}

			price := e.ChildText("span.fs-price-lockup__dollars") + "." + e.ChildText("span.fs-price-lockup__cents")
			unit := e.ChildText("span.fs-price-lockup__per")
			imageStyle := e.ChildAttr("div.fs-product-image__inner", "style")
			image := regexp.MustCompile(`http.[^)]+`).FindString(imageStyle)

			url := e.Request.URL.String()

			//fmt.Println(title + " > " + productId + " > " + price + "/" + unit + " ---> " + image)

			if productId != "" && price != "" {
				// 在缓存系统中校验是否已经保存过了当天的数据
				checkKey := time.Now().Format("20060102") + SPIDER_PAKNSAVE + productId
				if !cache.Has(checkKey) {

					cache.Set(checkKey, 1)
					var item model.Item
					if model.DB.Where("website = ? AND product_id = ?", SPIDER_PAKNSAVE, productId).First(&item).RecordNotFound() {
						// 没找到旧数据时，新建商品记录
						item.Image = image
						item.ProductID = productId
						item.InternalID = productId
						item.Title = title
						item.TitleZh = titleZh
						item.Website = SPIDER_PAKNSAVE
						item.Unit = unit
						item.Url = url
						model.DB.Create(&item)
					}

					flPrice, _ := strconv.ParseFloat(price, 10)
					model.DB.Model(&item).Association("Prices").Append(model.Price{Price: flPrice})
				}
			}
		})
		log.Println("PaknSave Run: " + w.url)
		w.cr.Visit(w.url)

	}
	return nil
}

func init() {
	// 在启动时注册Paknsave类工厂
	Register(SPIDER_PAKNSAVE, func() Spider {
		return new(Paknsave)
	})
}
