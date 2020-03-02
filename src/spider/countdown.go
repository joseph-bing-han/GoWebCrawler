package spider

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/chromed"
	"GoWebCrawler/src/utils/conf"
	"GoWebCrawler/src/utils/mq"
	gjson "encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/bregydoc/gtranslate"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Countdown struct {
	cr       *colly.Collector
	url      string
	cookies  string
	isUpdate bool
	tries    int
}

func (w *Countdown) SetURL(url string, isUpdate bool) {
	if w.cr == nil {
		w.cr = NewCollector(false)
	}
	w.url = url
	w.isUpdate = isUpdate
}

func (w *Countdown) Run() error {

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
				if strings.Contains(url, "/shop/recipe") {
					return
				}
				if strings.Contains(url, "/shop/productdetails") {
					return
				}
				url = "https://shop.countdown.co.nz" + url
				checkKey := SPIDER_COUNTDOWN + url
				// todo: test
				if !cache.Has(checkKey) {
					cache.Set(checkKey, 1)
					log.Println("[INFO]", "["+SPIDER_COUNTDOWN+"]", "Get URL: "+url)
					mq.Add(map[string]interface{}{"url": url, "update": false})
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

							itemId := product["sku"].(string)

							priceNode := product["price"].(map[string]interface{})

							price, _ := priceNode["salePrice"].(gjson.Number).Float64()

							unit := product["unit"].(string)

							imageNode := product["images"].(map[string]interface{})

							image := imageNode["big"].(string)
							image = "https://shop.countdown.co.nz" + image

							url := "https://shop.countdown.co.nz/shop/productdetails?stockcode=" + itemId + "&name=" + productId
							//strPrice := fmt.Sprintf("%f", price)
							//fmt.Println(title + "(" + titleZh + ") > " + productId + " > " + strPrice + "/" + unit + " ---> " + image)

							category, err := cache.Get("Category-" + url)
							if err != nil {
								category = ""
							}

							// 在缓存系统中校验是否已经保存过了当天的数据
							checkKey := SPIDER_COUNTDOWN + productId
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
									item.Url = url
									item.Category = category.(string)
									model.DB.Create(&item)
								} else {
									item.Url = url
									item.Category = category.(string)
									model.DB.Save(&item)
								}
								model.DB.Model(&item).Association("Prices").Append(model.Price{Price: price})
							}
						}

					}

				}
			}

		})

		w.cr.OnHTML("div#content-panel", func(e *colly.HTMLElement) {

			title := e.ChildText("h1")
			if title == "" {
				return
			}

			products := regexp.MustCompile(`name=([^&]+)`).FindStringSubmatch(e.Request.URL.String())
			productId := products[1]
			if title != "" && productId != "" {
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

				unit := ""
				itemId := e.ChildAttr("input[name='stockcode']", "value")
				priceText := e.ChildText("span.price")
				prices := regexp.MustCompile(`\$(\d+\.\d+)(.*)`).FindStringSubmatch(priceText)
				price := prices[1]
				if len(prices) == 3 {
					unit = strings.TrimSpace(prices[2])
				}

				image := e.ChildAttr("img.product-image", "src")
				image = "https://shop.countdown.co.nz" + image

				url := e.Request.URL.String()

				category, err := cache.Get("Category-" + url)
				if err != nil {
					category = ""
				}

				//fmt.Println(title + "(" + titleZh + ") " + category.(string) + " > " + productId + "[" + itemId + "] > " + price + "/" + unit + " ---> " + image)

				// 在缓存系统中校验是否已经保存过了当天的数据
				checkKey := SPIDER_COUNTDOWN + productId
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
						item.Url = url
						item.Category = category.(string)
						model.DB.Create(&item)
					} else {
						item.Image = image
						item.Unit = unit
						item.ProductID = productId
						item.InternalID = itemId
						item.Title = title
						item.TitleZh = titleZh
						item.Website = SPIDER_COUNTDOWN
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
			log.Println("[ERROR]", "["+SPIDER_COUNTDOWN+"]", err)
			w.tries--
			if w.tries >= 0 {
				time.Sleep(time.Second)
				response.Request.Retry()
			}

		})

		log.Println("[INFO]", "["+SPIDER_COUNTDOWN+"]", "RUN: "+w.url)
		w.cr.Visit(w.url)
	}
	return nil
}

func init() {
	var cookies string
	key := time.Now().Format("20060102") + SPIDER_COUNTDOWN + "-cookie-key"
	value, error := cache.Get(key)
	if error == nil || value.(string) != "" {
		cookies = value.(string)
	} else {
		cookies = chromed.GetCookie("https://shop.countdown.co.nz", key)
		cache.Set(key, cookies)
	}

	// 在启动时注册Countdown类工厂
	Register(SPIDER_COUNTDOWN, func() Spider {
		countdown := new(Countdown)
		countdown.cookies = cookies
		tries, err := strconv.Atoi(conf.Get("TRIES", "3"))
		if err != nil {
			tries = 3
		}
		countdown.tries = tries

		return countdown
	})
}
