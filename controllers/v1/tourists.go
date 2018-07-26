package v1

import (
	"estate/models/v1"
	"strconv"

	"github.com/gin-gonic/gin"
)

var touristsModel = new(v1.TouristsModel)

// 游客-路由
func Tourists(parentRoute *gin.RouterGroup) {
	router := parentRoute.Group("/tourists")
	router.GET("estate_detail", Tourists_EstateDetail)          // 4.2 游客-房源详情
	router.POST("estate_consulting", Tourists_EstateConsulting) // 4.3 游客-房源咨询
}

// 游客-房源列表
// func Tourists_EstastList(c *gin.Context) {
// 	estateCode := c.Query("estate_code")
// 	listorder, _ := strconv.Atoi(c.Query("listorder"))
// 	screenJson := c.Query("screen_json")
// 	perPage, _ := strconv.Atoi(c.Query("per_page"))
// 	lastId, _ := strconv.Atoi(c.Query("last_id"))

// 	// 获取房源列表
// 	data, errMsg := touristsModel.Tourists_EstastList(estateCode, screenJson, listorder, perPage, lastId)
// 	if errMsg != "" {
// 		c.JSON(400, gin.H{
// 			"code": 1010,
// 			"msg":  errMsg,
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"code": 0,
// 		"msg":  "success",
// 		"data": data,
// 	})
// 	return
// }

// 游客-房源信息
func Tourists_EstateDetail(c *gin.Context) {
	estateId, _ := strconv.Atoi(c.Query("estate_id"))
	if estateId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 获取房源信息
	data, errMsg := touristsModel.Tourists_EstateDetail(estateId)
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
	return
}

// 游客-房源咨询
func Tourists_EstateConsulting(c *gin.Context) {
	estateId, _ := strconv.Atoi(c.PostForm("estate_id"))
	name := c.PostForm("name")
	wechat := c.PostForm("wechat")
	sex, _ := strconv.Atoi(c.PostForm("sex"))
	if estateId == 0 || name == "" || wechat == "" || sex == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 提交咨询信息
	errMsg := touristsModel.Tourists_EstateConsulting(estateId, sex, name, wechat)
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}
	c.JSON(201, gin.H{
		"code": 0,
		"msg":  "success",
	})
	return
}
