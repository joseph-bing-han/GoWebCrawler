package main

import (
	"GoWebCrawler/src/utils/conf"
	"GoWebCrawler/src/web/router"
	"context"
	"github.com/foolin/goview"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
	"github.com/foolin/goview/supports/ginview"
)

func main() {

	/////////////////////////////////////////////////////////
	// 服务器模式相关
	// 设置服务器模式，默认是发布模式
	gin.SetMode(conf.Get("WEB_SERVER_MODE", "release"))
	/////////////////////////////////////////////////////////

	/////////////////////////////////////////////////////////
	// 日志相关
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()

	// 创建文件夹
	logFile := conf.Get("WEB_SERVER_LOG_FILE", "./logs/gin.log")
	path, _ := filepath.Split(logFile)
	_ = os.MkdirAll(path, os.ModePerm)

	// 记录到文件。
	f, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(f)

	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	////////////////////////////////////////////////////////

	gRouter := gin.Default()

	// 注册模板
	gRouter.HTMLRender = ginview.New(goview.Config{
		Root:         "resources/tpl",
		Extension:    ".gohtml",
		Master:       "layouts/base",
		DisableCache: conf.Get("WEB_SERVER_MODE", "release") == "debug",
	})

	// 注册静态资源
	gRouter.Static("/img", "resources/img")
	gRouter.Static("/js", "resources/js")
	gRouter.Static("/css", "resources/css")

	// 注册路由列表
	new(router.DefaultRouter).Register(gRouter)

	ip := conf.Get("WEB_SERVER_HOST", "127.0.0.1")
	port := conf.Get("WEB_SERVER_PORT", "8080")
	srv := &http.Server{
		Addr:    ip + ":" + port,
		Handler: gRouter,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
