package spider

import (
	"GoWebCrawler/src/model"
	"GoWebCrawler/src/utils/cache"
	"GoWebCrawler/src/utils/mq"
	"context"
	"github.com/bregydoc/gtranslate"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Kmart struct {
	cr      *colly.Collector
	url     string
	cookies string
}

func (k *Kmart) SetCookies(cookies string) {
	k.cookies = cookies
}

func (w *Kmart) SetURL(url string) {
	if w.cr == nil {
		w.cr = NewCollector()
	}
	w.url = url
}

func (w *Kmart) Run() error {

	if len(w.url) > 0 {
		w.cr.OnRequest(func(request *colly.Request) {
			request.Headers.Set("cookie", w.cookies)
		})

		// 处理所有链接
		w.cr.OnHTML("a[href]", func(e *colly.HTMLElement) {
			//fmt.Println(e)
			url := e.Attr("href")
			//fmt.Println("Get URL:" + url)
			if match, _ := regexp.MatchString(`^/[\w\W]+$`, url); match {
				url = "https://www.kmart.co.nz" + url
				checkKey :=  SPIDER_KMART + url
				// todo: test
				if !cache.Has(checkKey) {
					cache.Set(checkKey, 1)
					//log.Println("Add URL: " + url)
					mq.Add(map[string]interface{}{"url": url})
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

			titleZh, error := gtranslate.TranslateWithParams(
				title,
				gtranslate.TranslationParams{
					From:  "en",
					To:    "zh",
					Delay: time.Second * 2,
				},
			)
			if error != nil {
				titleZh = title
			}

			itemId := e.ChildText("h7")
			price := e.ChildAttr("span.price", "aria-label")
			price = strings.Replace(price, "$ ", "", -1)
			image := e.ChildAttr("img.mainImg", "src")
			image = "https://www.kmart.co.nz" + image

			url := e.Request.URL.String()

			//fmt.Println(title + "(" + titleZh + ") > " + productId + " > " + price + " ---> " + image)
			if productId != "" && price != "" {
				// 在缓存系统中校验是否已经保存过了当天的数据
				checkKey :=  SPIDER_KMART + productId
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
						model.DB.Create(&item)
					}

					flPrice, _ := strconv.ParseFloat(price, 10)
					model.DB.Model(&item).Association("Prices").Append(model.Price{Price: flPrice})
				}
			}
		})
		log.Println("Kmart Run: " + w.url)
		w.cr.Visit(w.url)

	}
	return nil
}

func getCookies() string {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var result string
	// navigate to a page, wait for an element, click
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.kmart.co.nz/`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > div#page`),
		// find and click "Expand All" link

		chromedp.ActionFunc(func(ctx context.Context) error {
			// 获取cookie
			cookies, err := network.GetAllCookies().Do(ctx)
			// 将cookie拼接成header请求中cookie字段的模式
			for _, v := range cookies {
				result = result + v.Name + "=" + v.Value + "; "

			}
			if err != nil {
				return err
			}
			return nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}
	return result
}

func init() {
	var cookies string
	key := time.Now().Format("20060102") + SPIDER_KMART + "-cookie-key"
	value, error := cache.Get(key)
	if error == nil || value.(string) != "" {
		cookies = value.(string)
	} else {
		cookies = getCookies()
		cache.Set(key, cookies)
	}
	//log.Println(cookies)
	// 在启动时注册Kmart类工厂
	Register(SPIDER_KMART, func() Spider {
		kmart := new(Kmart)
		kmart.SetCookies(cookies)
		return kmart
	})
}
