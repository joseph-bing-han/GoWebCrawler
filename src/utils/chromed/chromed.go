package chromed

import (
	"GoWebCrawler/src/utils/cache"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"time"
	"context"
)

func GetCookie(url string, cacheKey string) string {
	checkKey := "CHROME-RUN-" + url
	// todo: test
	if !cache.Has(checkKey) {
		cache.Set(checkKey, url)
		// create chrome instance
		ctx, cancel := chromedp.NewContext(
			context.Background(),
			chromedp.WithLogf(log.Printf),
		)
		defer cancel()

		// create a timeout
		ctx, cancel = context.WithTimeout(ctx, 3*time.Minute)
		defer cancel()

		var result string
		// navigate to a page, wait for an element, click
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),

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
			log.Println("[ERROR]", "[CHROMEDP]", err)
			result = ""
		}
		cache.Delete(checkKey)
		return result
	} else {
		for {
			time.Sleep(time.Second * 3)
			if !cache.Has(checkKey) {
				var result string
				value, error := cache.Get(cacheKey)
				if error == nil || value.(string) != "" {
					result = value.(string)
				}
				return result
			}
		}

	}

}
