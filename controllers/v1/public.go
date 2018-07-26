package v1

import (
	"estate/middleware"
	"estate/models/v1"
	"strconv"

	"github.com/gin-gonic/gin"
)

var publicModel = new(v1.PublicModel)

// 公用-路由
func Public(parentRoute *gin.RouterGroup) {
	router := parentRoute.Group("/public")
	router.Use(middleware.Auth())
	router.GET("/company_detail", Public_CompanyDetail)          // 3.1 公用-公司详情
	router.GET("/sales_manage/list", Public_SalesManageList)     // 3.2.1 公用-销售管理-列表
	router.GET("/sales_manage/detail", Public_SalesManageDetail) // 3.2.1 公用-销售管理-详情
}

// 公用-公司详情
func Public_CompanyDetail(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	if userId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 获取公司详情
	data, errMsg := publicModel.Public_CompanyDetail(userId)
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

// 公用-销售管理列表
func Public_SalesManageList(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	userType, _ := strconv.Atoi(c.Query("user_type"))
	if userId == 0 || userType == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if userType == 0 { // 非主管不允许看销售列表
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "该账户非公司主管，不允许看销售列表",
		})
		return
	}

	// 销售列表
	var data interface{}
	data, errMsg := publicModel.Public_SalesManageList(userId)
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}
	if data == nil {
		data = make(map[string]interface{})
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
	return
}

// 公用-销售管理详情
func Public_SalesManageDetail(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	if userId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 销售详情
	data, errMsg := publicModel.Public_SalesManageDetail(userId)
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
