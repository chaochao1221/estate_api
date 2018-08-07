package v1

import (
	"estate/models/v1"
	"estate/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var touristsModel = new(v1.TouristsModel)

// 游客-路由
func Tourists(parentRoute *gin.RouterGroup) {
	router := parentRoute.Group("/tourists")
	router.GET("estate_list", Tourists_EstateList)              // 4.1 游客-房源列表
	router.POST("estate_consulting", Tourists_EstateConsulting) // 4.2 游客-房源咨询
}

// 游客-房源列表
func Tourists_EstateList(c *gin.Context) {
	estParam := &(v1.PublicEstateListParamter{
		Keyword:    c.Query("keyword"),
		Listorder:  utils.Str2int(c.Query("listorder")),
		ScreenJson: c.Query("screen_json"),
		PerPage:    utils.Str2int(c.Query("per_page")),
		LastId:     utils.Str2int(c.Query("last_id")),
	})

	// 房源列表
	data, errMsg := touristsModel.Tourists_EstateList(estParam)
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}
	if data == nil {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": make(map[string]interface{}),
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
