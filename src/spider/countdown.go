package spider

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/mq"
	gjson "encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/bregydoc/gtranslate"
	"github.com/gocolly/colly"
	"log"
	"regexp"
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

		w.cr.OnHTML("div#product-list", func(e *colly.HTMLElement) {
			js := e.Text
			matches := regexp.MustCompile(`PRODUCT_GRI.*`).FindAllString(js, 1)
			if len(matches) == 1 {
				matches = regexp.MustCompile(`\[\{.*\}\]`).FindAllString(matches[0], 1)
				if len(matches) == 1 {
					json, _ := simplejson.NewJson([]byte(matches[0])) //反序列化
					nodes, _ := json.Array()
					for _, node := range nodes {
						product := node.(map[string]interface{})

						title := product["name"].(string)

						productId := product["slug"].(string)

						if title != "" && productId != "" {
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

							itemId := product["sku"].(string)

							priceNode := product["price"].(map[string]interface{})

							price, _ := priceNode["salePrice"].(gjson.Number).Float64()

							unit := product["unit"].(string)

							imageNode := product["images"].(map[string]interface{})

							image := imageNode["big"].(string)
							image = "https://shop.countdown.co.nz" + image

							//strPrice := fmt.Sprintf("%f", price)
							//fmt.Println(title + "(" + titleZh + ") > " + productId + " > " + strPrice + "/" + unit + " ---> " + image)

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
									item.InternalID = itemId
									item.Title = title
									item.TitleZh = titleZh
									item.Website = SPIDER_COUNTDOWN
									model.DB.Create(&item)
								}
								model.DB.Model(&item).Association("Prices").Append(model.Price{Price: price})
							}
						}

					}

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
