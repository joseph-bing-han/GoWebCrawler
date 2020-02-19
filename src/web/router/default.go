package router

import (
	"GoWebCrawler/src/web/controller"
	"github.com/gin-gonic/gin"
)

type DefaultRouter struct {
}

func (*DefaultRouter) Register(g *gin.Engine) {
	g.GET("/", controller.HomeIndex)
}
