package spider

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/chromed"
	"GoWebCrawler/src/utils/conf"
	"GoWebCrawler/src/utils/mq"
	"github.com/bregydoc/gtranslate"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Kmart struct {
	cr       *colly.Collector
	url      string
	cookies  string
	isUpdate bool
	tries    int
}

func (w *Kmart) SetURL(url string, isUpdate bool) {
	if w.cr == nil {
		w.cr = NewCollector(true)
	}
	w.url = url
	w.isUpdate = isUpdate
}

func (w *Kmart) Run() error {

	if len(w.url) > 0 {
		w.cr.OnRequest(func(request *colly.Request) {
			request.Headers.Set("cookie", w.cookies)
		})

		// 处理所有链接
		w.cr.OnHTML("a[href]", func(e *colly.HTMLElement) {
			if w.isUpdate {
				return
			}
			url := e.Attr("href")
			//fmt.Println("Get URL:" + url)
			if match, _ := regexp.MatchString(`^/[\w\W]+$`, url); match {
				url = "https://www.kmart.co.nz" + url
				checkKey := SPIDER_KMART + url
				// todo: test
				if !cache.Has(checkKey) {
					cache.Set(checkKey, 1)
					log.Println("[INFO]", "["+SPIDER_KMART+"]", "Get URL: "+url)
					mq.Add(map[string]interface{}{"url": url, "update": false})
				}
			}
		})

		// 处理商品页面数据
		w.cr.OnHTML("body", func(e *colly.HTMLElement) {
			productId := e.ChildAttr("input#productID", "value")

			if productId == "" {
				return
			}

			title := e.ChildText("h1")
			if title == "" {
				return
			}

			titleZh := title
			if !w.isUpdate {
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
			}

			itemId := e.ChildText("h7")
			price := e.ChildAttr("span.price", "aria-label")
			price = strings.Replace(price, "$ ", "", -1)
			image := e.ChildAttr("img.mainImg", "src")
			image = "https://www.kmart.co.nz" + image

			url := e.Request.URL.String()

			category, err := cache.Get("Category-" + url)
			if err != nil {
				category = ""
			}

			//fmt.Println(title + "(" + titleZh + ") > " + productId + " > " + price + " ---> " + image)
			if productId != "" && price != "" {
				// 在缓存系统中校验是否已经保存过了当天的数据
				checkKey := SPIDER_KMART + productId
				if !cache.Has(checkKey) {

					cache.Set(checkKey, 1)
					var item model.Item
					if model.DB.Where("website = ? AND product_id = ?", SPIDER_KMART, productId).First(&item).RecordNotFound() {
						// 没找到旧数据时，新建商品记录
						item.Image = image
						item.ProductID = productId
						item.InternalID = itemId
						item.Title = title
						item.TitleZh = titleZh
						item.Website = SPIDER_KMART
						item.Url = url
						item.Category = category.(string)
						model.DB.Create(&item)
					} else {
						item.Image = image
						item.ProductID = productId
						item.InternalID = itemId
						item.Title = title
						item.TitleZh = titleZh
						item.Website = SPIDER_KMART
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
			log.Println("[ERROR]", "["+SPIDER_KMART+"]", err)

			w.tries--
			if w.tries >= 0 {
				time.Sleep(time.Second)
				response.Request.Retry()
			}

		})

		log.Println("[INFO]", "["+SPIDER_KMART+"]", "RUN: "+w.url)
		w.cr.Visit(w.url)

	}
	return nil
}

func init() {
	var cookies string
	key := time.Now().Format("20060102") + SPIDER_KMART + "-cookie-key"
	value, error := cache.Get(key)
	if error == nil || value.(string) != "" {
		cookies = value.(string)
	} else {
		cookies = chromed.GetCookie("https://www.kmart.co.nz/", key)
		cache.Set(key, cookies)
	}
	//log.Println(cookies)
	// 在启动时注册Kmart类工厂
	Register(SPIDER_KMART, func() Spider {
		kmart := new(Kmart)
		kmart.cookies = cookies
		tries, err := strconv.Atoi(conf.Get("TRIES", "3"))
		if err != nil {
			tries = 3
		}
		kmart.tries = tries
		return kmart
	})
}
