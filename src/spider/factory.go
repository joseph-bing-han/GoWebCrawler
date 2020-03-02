package spider

import (
	"GoWebCrawler/src/utils/conf"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"log"
	"strings"
	"time"
)

type Spider interface {
	SetURL(url string, isUpdate bool)
	Run() error
}

var (
	// 保存注册好的工厂信息
	factoryByName = make(map[string]func() Spider)
)

// 注册一个类生成工厂
func Register(name string, factory func() Spider) {
	factoryByName[name] = factory
}

// 根据名称创建对应的类
func Create(name string) Spider {
	if f, ok := factoryByName[name]; ok {
		return f()
	} else {
		panic("name not found")
	}
}

func NewCollector(defaultProxy bool) *colly.Collector {
	cr := colly.NewCollector()

	// 使用随机User Agent
	extensions.RandomUserAgent(cr)

	// 添加referer防盗链
	extensions.Referer(cr)

	// 添加必要的浏览器信息
	cr.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Connection", "keep-alive")
		request.Headers.Set("Cache-Control", "'max-age=0")
		request.Headers.Set("Upgrade-Insecure-Request", "1")
		request.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		request.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,ja;q=0.8,en-NZ;q=0.7,en;q=0.6")

	})

	// 初始化代理池
	var proxyIP []string

	if defaultProxy {
		proxyIP = strings.Split(conf.Get("TOR_PROXY", "socks5://xebni:xebni@13.239.73.54:1984"), ",")
	} else {
		proxyIP = strings.Split(conf.Get("ALT_PROXY", "socks5://xebni:xebni@13.239.73.54:1984"), ",")
	}

	ps, err := proxy.RoundRobinProxySwitcher(proxyIP...)

	if err != nil {
		log.Fatalln(err)
	}
	cr.SetProxyFunc(ps)

	// 加了代理池，速度变慢，超时延长为2分钟
	cr.SetRequestTimeout(time.Minute * 2)
	return cr
}
