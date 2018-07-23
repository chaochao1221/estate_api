package main

import (
	"estate/db"
	"estate/pkg"
	"estate/routes"
	"flag"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	httpport := flag.Int("p", 8001, "HttpPort")
	flag.Parse()
	// 初始化...
	pkg.Init()
	// 全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	gin.SetMode(gin.DebugMode)
	// Initialize db
	db.Init()
	defer db.Close()
	// Initialize the routes
	router := routes.Init()
	// listen and serve on 0.0.0.0:8010
	router.Run(":" + strconv.Itoa(*httpport))
}
