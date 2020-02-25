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

var NEWWORLD_STORES = [...]string{
	"browser_nearest_store={\"UserLat\":\"-37.667256\",\"UserLng\":\"175.150107\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.905625\",\"UserLng\":\"174.994131\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.02246\",\"UserLng\":\"174.89719\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.623065\",\"UserLng\":\"175.282791\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.728271\",\"UserLng\":\"175.273985\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.359638\",\"UserLng\":\"175.600741\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.117024\",\"UserLng\":\"174.893023\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.948688\",\"UserLng\":\"175.665517\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.763486\",\"UserLng\":\"175.153556\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.753024\",\"UserLng\":\"175.239998\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.774933\",\"UserLng\":\"174.553672\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.279416\",\"UserLng\":\"174.780508\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.064268\",\"UserLng\":\"174.941101\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.671351\",\"UserLng\":\"176.877001\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.628789\",\"UserLng\":\"174.732596\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.890274\",\"UserLng\":\"174.831971\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.201793\",\"UserLng\":\"175.867411\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.926194\",\"UserLng\":\"175.041161\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.834308\",\"UserLng\":\"175.699435\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-38.146066\",\"UserLng\":\"176.237248\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.158214\",\"UserLng\":\"174.20708\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.314466\",\"UserLng\":\"174.780574\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.822527\",\"UserLng\":\"175.287416\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.18988\",\"UserLng\":\"174.90306\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.715961\",\"UserLng\":\"174.747287\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.05782\",\"UserLng\":\"174.07774\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.82951\",\"UserLng\":\"174.796193\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.801952\",\"UserLng\":\"175.322549\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.348551\",\"UserLng\":\"175.623806\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.99467\",\"UserLng\":\"176.55584\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.910865\",\"UserLng\":\"174.685569\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.387171\",\"UserLng\":\"175.639157\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.684964\",\"UserLng\":\"174.73973\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.209025\",\"UserLng\":\"174.90804\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-35.716277\",\"UserLng\":\"174.321962\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.933874\",\"UserLng\":\"174.911485\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-35.2267\",\"UserLng\":\"173.951368\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.033816\",\"UserLng\":\"177.416989\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.000167\",\"UserLng\":\"174.236804\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.893425\",\"UserLng\":\"175.472202\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.811428\",\"UserLng\":\"174.711486\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.210911\",\"UserLng\":\"176.096788\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.644715\",\"UserLng\":\"176.847154\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-38.085853\",\"UserLng\":\"176.700944\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.903857\",\"UserLng\":\"174.925067\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.07666\",\"UserLng\":\"175.379646\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.452412\",\"UserLng\":\"175.841836\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.472517\",\"UserLng\":\"175.281707\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.9802\",\"UserLng\":\"174.853287\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-38.218917\",\"UserLng\":\"175.867537\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.315835\",\"UserLng\":\"174.814635\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.874943\",\"UserLng\":\"175.066744\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.713787\",\"UserLng\":\"176.142697\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.249975\",\"UserLng\":\"174.727925\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.728207\",\"UserLng\":\"174.710519\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.026035\",\"UserLng\":\"175.526349\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.772298\",\"UserLng\":\"174.764805\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-38.011812\",\"UserLng\":\"177.275962\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-38.991014\",\"UserLng\":\"175.809118\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-38.331839\",\"UserLng\":\"175.16319\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.846405\",\"UserLng\":\"174.765935\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.783124\",\"UserLng\":\"176.328028\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.803437\",\"UserLng\":\"175.767623\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.284257\",\"UserLng\":\"174.738052\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.508662\",\"UserLng\":\"176.887458\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-35.408951\",\"UserLng\":\"173.801012\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.950929\",\"UserLng\":\"176.993759\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.292409\",\"UserLng\":\"174.784342\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.203041\",\"UserLng\":\"174.807839\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.287948\",\"UserLng\":\"174.775268\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.148436\",\"UserLng\":\"175.012536\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.393116\",\"UserLng\":\"175.839644\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.136386\",\"UserLng\":\"174.841919\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.586612\",\"UserLng\":\"174.69366\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.860989\",\"UserLng\":\"174.829234\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.223472\",\"UserLng\":\"174.823202\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.659507\",\"UserLng\":\"175.522299\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.337111\",\"UserLng\":\"174.286329\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-40.228015\",\"UserLng\":\"175.569532\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-38.882994\",\"UserLng\":\"175.257912\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.274006\",\"UserLng\":\"174.778067\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.247527\",\"UserLng\":\"174.791469\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.06659\",\"UserLng\":\"174.102988\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-39.525324\",\"UserLng\":\"176.862409\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.174103\",\"UserLng\":\"174.980844\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.848637\",\"UserLng\":\"174.751367\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-35.755737\",\"UserLng\":\"174.367937\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.908622\",\"UserLng\":\"174.734362\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.169962\",\"UserLng\":\"174.82611\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.334534\",\"UserLng\":\"174.772497\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-36.398709\",\"UserLng\":\"174.666795\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.653862\",\"UserLng\":\"176.198938\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-41.094112\",\"UserLng\":\"174.868615\",\"IsSuccess\":true}",
	"browser_nearest_store={\"UserLat\":\"-37.688946\",\"UserLng\":\"176.1347\",\"IsSuccess\":true}",
}

var NEWWORLD_BRANCH = [...]string{
	"New World Ngaruawahia (7 Galileo Street, Ngaruawahia, 3720)",
	"New World Kapiti (159 Kapiti Road, Paraparaumu, 5032)",
	"New World Southmall (187 Great South Road, Manurewa, Auckland, 2102)",
	"New World Levin (21 Bath Street, Levin, 5510)",
	"New World Rototuna (Cnr Thomas & Horsham Downs Rds, Hamilton, 3210)",
	"New World Pioneer (179-197 Main Street, Palmerston North, 4412)",
	"New World Whitby (Discovery Drive, Whitby Village Centre, Wellington, 5024)",
	"New World Masterton (Cnr Bruce and Dixon Street, Masterton, 5810)",
	"New World Otaki (155-163 Main Highway, Otaki, 5512)",
	"New World Te Rapa (751 Te Rapa Road, Te Rapa, Hamilton, 3200)",
	"New World Kumeu (110 Main Road, Auckland, 0810)",
	"New World Railway Station (2 Bunny Street, Pipitea, Wellington, 6011)",
	"New World Papakura (29-31 East Street, Auckland, 2110)",
	"New World Havelock North (Porter Drive, Havelock North, 4130)",
	"New World Whangaparaoa (588 Whangaparaoa Road, Stanmore Bay, Auckland, 0932)",
	"New World Stonefields (100 Lunn Avenue, Mt Wellington, Auckland, 1072)",
	"New World Whangamata (308 Acikin Road, Whangamata, 3620)",
	"New World Wanganui (374 Victoria Avenue, Wanganui, 4500)",
	"New World Whitianga (1 Joan Gaskell Drive, Whitianga, 3510)",
	"New World Westend (247 Old Taupo Road, Rotorua, 3015)",
	"New World Inglewood (46 Matai Street, Inglewood, 4330)",
	"New World Newtown (195 Riddiford Street, Wellington, 6021)",
	"New World Glenview (Ohaupo Road, Glenview, Hamilton, 3206)",
	"New World Pukekohe (17 Paerata Road, Pukekohe, Auckland, 2120)",
	"New World Browns Bay (2 Inverness Road, Auckland, 0630)",
	"New World New Plymouth (78 Courtenay Street, New Plymouth, 4310)",
	"New World Devonport (35 Bartley Terrace, Devonport, Auckland, 0624)",
	"New World Hillcrest (280 Cambridge Road, Hillcrest, Hamilton, 3216)",
	"New World Broadway (Broadway Avenue, Palmerston North, 4414)",
	"New World Waipukurau (27 Russell Street, Waipukurau, 4200)",
	"New World New Lynn (2-6 Crown Lynn Place, New Lynn, Auckland, 0600)",
	"New World Aokautere (194-200 Ruapehu Drive, Summerhill, Aokautere, 4410)",
	"New World Long Bay (55B Glenvar Ridge Road, Long Bay, Auckland, 0630)",
	"New World Hutt City (Bloomfield Terrace, Lower Hutt, 5010)",
	"New World Regent (167 Bank Street, Regent, Whangarei, 0112)",
	"New World Botany (588 Chapel Road, East Tamaki, Auckland, 2013)",
	"New World Kerikeri (99 Kerikeri Road, Kerikeri, 0230)",
	"New World Wairoa (41 Queen Street, Wairoa, 4108)",
	"New World Waitara (42 Queen Street, Waitara, 4320)",
	"New World Cambridge (14 Anzac Street, Cambridge, 3434)",
	"New World Birkenhead (180 Mokoia Road, Chatswood, Auckland, 0626)",
	"New World Dannevirke (Denmark Street, Dannevirke, 4930)",
	"New World Hastings (400 Heretaunga Street East, Hastings, 4122)",
	"New World Kawerau (Tarawera Court, Kawerau, 3127)",
	"New World Howick (77 Union Road, Howick, Auckland, 2014)",
	"New World Marton (427 Wellington Road, Marton, 4710)",
	"New World Pahiatua (101 Main Street, Pahiatua, 4910)",
	"New World Foxton (Cnr Main and Whyte Streets, FOXTON, 4814)",
	"New World Papatoetoe (65 St Georges Street, Papatoetoe, Auckland, 2025)",
	"New World Tokoroa (Bridge Street, Tokoroa, 3420)",
	"New World Miramar (48 Miramar Avenue, Miramar, WELLINGTON, 6022)",
	"New World Waikanae (5 Parata Street, Waikanae, 5036)",
	"New World Gate Pa (948 Cameron Road, Gate Pa, Tauranga, 3112)",
	"New World Waiuku (25-49 Bowen Street, Waiuku, 2123)",
	"New World Albany (219 Don McKinnon Drive, Albany, Auckland, 0632)",
	"New World Carterton (60 High Street South, Carterton, 5713)",
	"New World Milford (141 Kitchener Road, Milford, Auckland, 0620)",
	"New World Opotiki (19 Bridge Street, Opotiki, 3122)",
	"New World Turangi (19 Ohuanga Street, Turangi, 3334)",
	"New World Te Kuiti (39-51 Rora Street, Te Kuiti, 3910)",
	"New World Metro Auckland (125 Queen Street, Auckland, 1010)",
	"New World Te Puke (12 Jocelyn Street, Te Puke, 3119)",
	"New World Matamata (45 Waharoa East Road, Matamata, 3400)",
	"New World Karori (236 Karori Road, Karori, Wellington, 6012)",
	"New World Onekawa (34 Maadi Road, Onekawa, Napier, 4110)",
	"New World Kaikohe (Marino Place, Kaikohe, 0405)",
	"New World Whakatane (51 Kakahoroa Drive, Whakatane, 3120)",
	"New World Wellington City (279 Wakefield Street, Wellington, 6011)",
	"New World Churton Park (Cnr Westchester Dr and Lakewood Av, Churton Park, Wellington, 6037)",
	"New World Metro (70 Willis Street, Wellington, 6011)",
	"New World Silverstream (28 Whitemans Road, Silverstream, Upper Hutt, 5019)",
	"New World Waihi (35 Kenny Street, Waihi, 3610)",
	"New World Porirua (Lyttleton Avenue, Porirua, 5022)",
	"New World Orewa (11 Moana Avenue, Orewa, 0931)",
	"New World Eastridge (209 Kepa Road, Mission Bay, Auckland, 1071)",
	"New World Newlands (Cnr Bracken and Newlands Road, Newlands, Wellington, 6037)",
	"New World Morrinsville (89 Thames Street, Morrinsville, 3300)",
	"New World Stratford (124 Regan Street, Stratford, 4332)",
	"New World Feilding (42 Aorangi Street, Feilding, 4702)",
	"New World Taumarunui (10 Hakiaha Street, Taumarunui, 3920)",
	"New World Thorndon (41 Murphy Street, Thorndon, Wellington, 6001)",
	"New World Khandallah (26 Ganges Road, Khandallah, Wellington, 6035)",
	"New World Merrilands (200 Mangorei Road, Merrilands, New Plymouth, 4312)",
	"New World Greenmeadows (9 Gloucester Street, Greenmeadows, Napier, 4112)",
	"New World Stokes Valley (14 Oates Street, Stokes Valley, Lower Hutt, 5019)",
	"New World Victoria Park (2 College Hill, Freemans Bay, Auckland, 1011)",
	"New World Onerahi (128 Onerahi Road, Whangarei, 0110)",
	"New World Mt Roskill (53 May Road, Mount Roskill, Auckland, 1041)",
	"New World Tawa (37 Oxford Terrace, Tawa, Wellington, 5028)",
	"New World Island Bay (6 Medway Street, Island Bay, Wellington, 6023)",
	"New World Warkworth (6 Percy Street, Warkworth, 0910)",
	"New World Mount Maunganui (Cnr Tweed St & Maunganui Rd, Mount Maunganui, 3116)",
	"New World Paremata (93-97 Mana Esplanade, Paremata, Porirua, 5026)",
	"New World Brookfield (Bellevue Road, Otumoetai, Tauranga, 3110)",
}

type NewWorld struct {
	cr       *colly.Collector
	url      string
	branch   int
	lowPrice float64
	isUpdate bool
}

func (w *NewWorld) SetURL(url string,isUpdate bool) {
	if w.cr == nil {
		w.cr = NewCollector(false)
	}
	w.url = url
	w.isUpdate = isUpdate
}

func (w *NewWorld) Run() error {

	if len(w.url) > 0 {

		w.cr.OnRequest(func(request *colly.Request) {
			if w.branch == -1 {
				cookie := NEWWORLD_STORES[0]
				request.Headers.Set("cookie", cookie)
			}
		})
		//处理所有链接
		w.cr.OnHTML("a[href]", func(e *colly.HTMLElement) {
			if w.isUpdate {
				return
			}
			url := e.Attr("href")
			if match, _ := regexp.MatchString(`^/[\w\W]+$`, url); match {
				url = "https://www.ishopnewworld.co.nz" + url
				checkKey := SPIDER_NEWWORLD + url
				// todo: test
				if !cache.Has(checkKey) {
					cache.Set(checkKey, 1)
					log.Println("[INFO]", "["+SPIDER_NEWWORLD+"]", "Get URL: "+url)
					mq.Add(map[string]interface{}{"url": url, "update": false})
				}
			}
		})

		// 处理商品页面数据
		w.cr.OnHTML("section.fs-product-detail", func(e *colly.HTMLElement) {
			defer func() {
				if len(NEWWORLD_STORES) > w.branch && w.branch >= 0 {
					cookie := NEWWORLD_STORES[w.branch]
					w.branch++
					e.Request.Headers.Set("cookie", cookie)
					e.Request.Retry()
				}
			}()

			productId := ""
			optionsJson := e.ChildAttr("div.fs-product-detail__wishlist", "data-options")
			if len(optionsJson) > 0 {
				var options map[string]interface{}
				json.Unmarshal([]byte(optionsJson), &options)
				productId = options["productId"].(string)
			}
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

			price := e.ChildText("span.fs-price-lockup__dollars") + "." + e.ChildText("span.fs-price-lockup__cents")
			unit := e.ChildText("span.fs-price-lockup__per")
			imageStyle := e.ChildAttr("div.fs-product-image__inner", "style")
			image := regexp.MustCompile(`http.[^)]+`).FindString(imageStyle)

			url := e.Request.URL.String()

			flPrice, _ := strconv.ParseFloat(price, 10)
			if productId != "" && price != "" && flPrice <= w.lowPrice {
				//fmt.Println(w.branch, ">>>"+title+"("+titleZh+") > "+productId+" > "+price+"/"+unit+" ---> "+image)
				// 在缓存系统中校验是否已经保存过了当天的数据
				checkKey := SPIDER_NEWWORLD + productId
				if !cache.Has(checkKey) {

					cache.Set(checkKey, 1)
					var item model.Item
					if model.DB.Where("website = ? AND product_id = ?", SPIDER_NEWWORLD, productId).First(&item).RecordNotFound() {
						// 没找到旧数据时，新建商品记录
						item.Image = image
						item.ProductID = productId
						item.InternalID = productId
						item.Title = title
						item.TitleZh = titleZh
						item.Unit = unit
						item.Website = SPIDER_NEWWORLD
						item.Url = url
						model.DB.Create(&item)
					}

					branch := NEWWORLD_BRANCH[0]
					if w.branch > 0 {
						branch = NEWWORLD_BRANCH[w.branch-1]
					}

					model.DB.Model(&item).Association("Prices").Append(model.Price{Price: flPrice, Branch: branch})
				} else {
					var item model.Item
					if !model.DB.Where("website = ? AND product_id = ?", SPIDER_NEWWORLD, productId).First(&item).RecordNotFound() {
						// 找到旧数据时，更新商品价格记录
						var price model.Price
						if !model.DB.Where("item_id = ? AND created_at >= ?", item.ID, time.Now().Format("2006-01-01 00:00:00")).First(&price).RecordNotFound() {
							branch := NEWWORLD_BRANCH[0]
							if w.branch > 0 {
								branch = NEWWORLD_BRANCH[w.branch-1]
							}
							price.Price = flPrice
							price.Branch = branch
							model.DB.Save(&price)
						}
					}
				}

				if w.branch == -1 {
					w.branch++
				}
				w.lowPrice = flPrice
			}
		})

		w.cr.OnError(func(response *colly.Response, err error) {
			log.Println("[ERROR]", "["+SPIDER_NEWWORLD+"]", err)
			time.Sleep(time.Second)
			response.Request.Retry()
		})

		log.Println("[INFO]", "["+SPIDER_NEWWORLD+"]", "RUN: "+w.url)
		w.branch = -1
		w.lowPrice = 99999
		w.cr.Visit(w.url)

	}
	return nil
}

func init() {
	// 在启动时注册NewWorld类工厂
	Register(SPIDER_NEWWORLD, func() Spider {
		return new(NewWorld)
	})
}
