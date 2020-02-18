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

var STORES = [...]string{
	"browser_nearest_store={\"UserLat\":\"-41.224456\",\"UserLng\":\"174.872537\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.03242\",\"UserLng\":\"174.867278\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.620413\",\"UserLng\":\"174.672975\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-35.09945\",\"UserLng\":\"173.258322\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.91014\",\"UserLng\":\"174.77342\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.89305\",\"UserLng\":\"174.70624\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.795653\",\"UserLng\":\"175.282839\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.205157\",\"UserLng\":\"174.899612\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.131559\",\"UserLng\":\"174.841818\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.636804\",\"UserLng\":\"176.836635\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.359198\",\"UserLng\":\"175.611537\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.053393\",\"UserLng\":\"174.93096\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.70191\",\"UserLng\":\"176.283687\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-35.725958\",\"UserLng\":\"174.32475\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.913347\",\"UserLng\":\"174.84009\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.319324\",\"UserLng\":\"174.79672\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.820257\",\"UserLng\":\"174.608696\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.509905\",\"UserLng\":\"176.869414\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.968516\",\"UserLng\":\"174.796863\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.205578\",\"UserLng\":\"174.913149\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-38.685532\",\"UserLng\":\"176.073247\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.738653\",\"UserLng\":\"176.103745\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-38.008101\",\"UserLng\":\"175.340102\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.696716\",\"UserLng\":\"176.161533\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.961198\",\"UserLng\":\"176.98257\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.05712\",\"UserLng\":\"174.08127\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.92638\",\"UserLng\":\"175.038901\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.857794\",\"UserLng\":\"174.627968\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.494788\",\"UserLng\":\"176.913562\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.7799\",\"UserLng\":\"175.27281\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.123787\",\"UserLng\":\"175.066601\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.930652\",\"UserLng\":\"174.913004\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.875918\",\"UserLng\":\"174.855148\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.987929\",\"UserLng\":\"174.880604\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.138597\",\"UserLng\":\"175.538994\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.729975\",\"UserLng\":\"174.706725\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.957214\",\"UserLng\":\"175.649687\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.918086\",\"UserLng\":\"175.001761\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.590947\",\"UserLng\":\"174.283801\",\"IsSuccess\":true}",
}

var BRANCH = [...] string{
	"PAK'nSAVE Lower Hutt (Brunswick Street, Hutt Central, Lower Hutt, 5010)",
	"PAK'nSAVE Royal Oak (691 Manukau Road, Royal Oak, Auckland, 1023)",
	"PAK'nSAVE Te Awamutu (650-670 Cambridge Road, Te Awamutu, 3800)",
	"PAK'nSAVE Kaitaia (111 North Road, Kaitaia, 0482)",
	"PAK'nSAVE New Plymouth (53 Leach Street, New Plymouth, 4310)",
	"PAK'nSAVE Tauriko (2 Taurikura Drive, Tauriko, Tauranga, 3110)",
	"PAK'nSAVE Manukau (6 Cavendish Drive, Manukau City, Auckland, 2104)",
	"PAK'nSAVE Mill Street (17 Mill Street, Whitiora, Hamilton, 3200)",
	"PAK'nSAVE Papakura (331-345 Great South Road, Takanini, Auckland, 2110)",
	"PAK'nSAVE Silverdale (20 Hibiscus Coast Highway, Silverdale, Auckland, 0932)",
	"PAK'nSAVE Wanganui (167 Glasgow Street, Wanganui, 4500)",
	"PAK'nSAVE Porirua (12 Parumoana Street, Porirua, 5022)",
	"PAK'nSAVE Napier City (25 Munroe Street, Napier South, Napier South NAPIER, 4110)",
	"PAK'nSAVE Whangarei (Walton Street, Whangarei, 0110)",
	"PAK'nSAVE Thames (100 Mary Street, Thames, 3500)",
	"PAK'nSAVE Cameron Road (476 Cameron Road, Tauranga, 3110)",
	"PAK'nSAVE Hawera (54 Princes Street, Hawera, 4610)",
	"PAK'nSAVE Kilbirnie (76 Rongotai Road, Kilbirnie, Wellington, 6003)",
	"PAK'nSAVE Clarence St (85 Clarence Street, Hamilton Lake, Hamilton, 3204)",
	"PAK'nSAVE Clendon (16 Robert Ross Place, Clendon Park, Auckland, 2103)",
	"PAK'nSAVE Tamatea (Leicester Avenue, Tamatea, Napier, 4112)",
	"PAK'nSAVE Upper Hutt (Gibbons Street, Upper Hutt, 5018)",
	"PAK'nSAVE Mt Albert (1167-1177 New North Road, Auckland, 1025)",
	"PAK'nSAVE Glen Innes (182 Apirana Avenue, Glen Innes, Auckland, 1072)",
	"PAK'nSAVE Mangere (44 Orly Avenue, Mangere, Auckland, 2022)",
	"PAK'nSAVE Masterton (Queen Street, Kuripuni, Masterton, 5810)",
	"PAK'nSAVE Westgate (17-19 Fred Taylor Drive, Massey, Auckland, 0814)",
	"PAK'nSAVE Lincoln Road (202 Lincoln Road, Henderson, Auckland, 0610)",
	"PAK'nSAVE Kapiti (76 Rimu Road, Paraparaumu, 5032)",
	"PAK'nSAVE Taupo (105-131 Ruapehu Street, Taupo, 3330)",
	"PAK'nSAVE Palmerston N (Fergusson Street, Palmerston North, 4472)",
	"PAK'nSAVE Petone (114-124 Jackson Street, Petone, Wellington, 5012)",
	"PAK'nSAVE Hastings (Heretaunga Street West, West End Shopping Centre, Hastings, 4120)",
	"PAK'nSAVE Pukekohe (99 Queen Street, Pukekohe, 2120)",
	"PAK'nSAVE Sylvia Park (286 Mt Wellington Highway, Mt Wellington, Auckland, 1060)",
	"PAK'nSAVE Whakatane (45 King Street, Whakatane, 3120)",
	"PAK'nSAVE Botany (501 Ti Rakau Drive, Northpark, Auckland, 2013)",
	"PAK'nSAVE Papamoa (42 Domain Road, Papamoa Beach, Papamoa, 3118)",
	"PAK'nSAVE Albany (Don McKinnon Drive, Albany, Auckland, 0632)",
}

type Paknsave struct {
	cr       *colly.Collector
	url      string
	branch   int
	lowPrice float64
}

func (w *Paknsave) SetURL(url string) {
	if w.cr == nil {
		w.cr = NewCollector(false)
	}
	w.url = url
}

func (w *Paknsave) Run() error {

	if len(w.url) > 0 {

		w.cr.OnRequest(func(request *colly.Request) {
			if w.branch == -1 {
				cookie := STORES[0]
				request.Headers.Set("cookie", cookie)
			}
		})

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
		w.cr.OnHTML(".fs-product-detail,.js-breadcrumbs", func(e *colly.HTMLElement) {

			defer func() {
				if len(STORES) > w.branch && w.branch >= 0 {
					cookie := STORES[w.branch]
					w.branch++
					e.Request.Headers.Set("cookie", cookie)
					e.Request.Retry()
				}
			}()

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
			//fmt.Println(title + "(" + titleZh + ") > " + productId + " > " + price + "/" + unit + " ---> " + image)
			flPrice, _ := strconv.ParseFloat(price, 10)
			if productId != "" && price != "" && flPrice <= w.lowPrice {
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

					branch := BRANCH[0]
					if w.branch >= 0 {
						branch = BRANCH[w.branch]
					}

					model.DB.Model(&item).Association("Prices").Append(model.Price{Price: flPrice, Branch: branch})
				}

				if w.branch == -1 {
					w.branch++
				}
				w.lowPrice = flPrice
			}
		})

		log.Println("PaknSave Run: " + w.url)
		w.branch = -1
		w.lowPrice = 99999
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
