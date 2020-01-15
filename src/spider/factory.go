package spider

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"log"
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

func NewCollector() *colly.Collector {
	cr := colly.NewCollector()

	// 使用随机User Agent
	extensions.RandomUserAgent(cr)

	// 添加referer防盗链
	extensions.Referer(cr)

	// 初始化代理池
	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:9010", "socks5://127.0.0.1:9020",
		"socks5://127.0.0.1:9030", "socks5://127.0.0.1:9040", "socks5://127.0.0.1:9050")
	if err != nil {
		log.Fatalln(err)
	}
	cr.SetProxyFunc(rp)

	// 加了代理池，速度变慢，超时延长为1分钟
	cr.SetRequestTimeout(time.Minute)
	return cr
}
