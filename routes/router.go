// @APIVersion 1.0.0
// @Apititle K12统考阅卷 API
// @Apidescription <a href="http://192.168.122.246:9004" target="_blank">k12统考阅卷文档</a>
package routes

import (
	"estate/controllers"
	// "estate/controllers/v1"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.Default()
	pprof.Register(router, nil)
	router.GET("/", controllers.Default)

	// api_v1 := router.Group("/v1")
	// {
	// 	v1.User(api_v1)
	// }

	// catch no router
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, "你这是闹哪样!")
	})
	return router
}
