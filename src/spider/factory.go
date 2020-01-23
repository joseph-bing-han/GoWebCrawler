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
	SetURL(url string)
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

	//// 初始化代理池
	var proxyIP [] string

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

	// 加了代理池，速度变慢，超时延长为1分钟
	cr.SetRequestTimeout(time.Minute * 2)
	return cr
}
