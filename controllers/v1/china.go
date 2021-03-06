package v1

import (
	"estate/middleware"
	"estate/models/v1"
	"estate/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var chinaModel = new(v1.ChinaModel)

// 中国中介-路由
func China(parentRoute *gin.RouterGroup) {
	router := parentRoute.Group("/china")
	router.Use(middleware.Auth())
	router.POST("/recommend", China_Recommend)
	router.GET("/customer_progress/list", China_CustomerProgressList)
	router.DELETE("/customer_progress/del/:id", China_CustomerProgressDel)
}

// 中国中介-推荐
func China_Recommend(c *gin.Context) {
	estateId, _ := strconv.Atoi(c.PostForm("estate_id"))
	name := c.PostForm("name")
	wechat := c.PostForm("wechat")
	sex, _ := strconv.Atoi(c.PostForm("sex"))
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if estateId == 0 || name == "" || wechat == "" || sex == 0 || userId == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 2 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "该用户非中国中介，不允许推荐客户",
		})
		return
	}

	// 推荐
	errMsg := chinaModel.China_Recommend(estateId, sex, userId, groupId, name, wechat)
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

// 中国中介-客户进展列表
func China_CustomerProgressList(c *gin.Context) {
	cusParam := &(v1.ChinaCustomerProgressListParameter{
		Keyword:     c.Query("keyword"),
		IsButt:      c.Query("is_butt"),
		IsToJapan:   c.Query("is_to_japan"),
		IsAgree:     c.Query("is_agree"),
		IsPay:       c.Query("is_pay"),
		IsLoan:      c.Query("is_loan"),
		UserId:      utils.Str2int(c.Query("user_id")),
		PerPage:     utils.Str2int(c.Query("per_page")),
		LastId:      utils.Str2int(c.Query("last_id")),
		LoginUserId: utils.Str2int(c.Request.Header.Get("user_id")),
		UserType:    utils.Str2int(c.Request.Header.Get("user_type")),
	})
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if cusParam.LoginUserId == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 2 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "该用户非中国中介，不允许查看客户进展",
		})
		return
	}

	// 客户进展
	data, errMsg := chinaModel.China_CustomerProgressList(cusParam)
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

// 中国中介-客户进展删除
func China_CustomerProgressDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if id == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 2 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "该用户非中国中介主管，不允许删除客户进展",
		})
		return
	}

	// 客户进展
	errMsg := chinaModel.China_CustomerProgressDel(id)
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
	})
	return
}
