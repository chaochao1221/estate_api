package v1

import (
	"estate/middleware"
	"estate/models/v1"
	"strconv"

	"github.com/gin-gonic/gin"
)

var japanModel = new(v1.JapanModel)

// 日本中介-路由
func Japan(parentRoute *gin.RouterGroup) {
	router := parentRoute.Group("/japan")
	router.Use(middleware.Auth())
	router.GET("/estate_progress", Japan_EstateProgress) // 6.1 日本中介-房源进展
}

// 日本中介-房源进展
func Japan_EstateProgress(c *gin.Context) {
	status, _ := strconv.Atoi(c.Query("status"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	lastId, _ := strconv.Atoi(c.Query("last_id"))
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if userId == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 3 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "该用户非日本中介，不允许查看房源进展",
		})
		return
	}

	// 房源进展
	data, errMsg := japanModel.Japan_EstateProgress(status, perPage, lastId, userId, userType)
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
