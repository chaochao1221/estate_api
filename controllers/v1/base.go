package v1

import (
	"estate/middleware"
	"estate/models/v1"
	"estate/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var baseModel = new(v1.BaseModel)

// 本部中介-路由
func Base(parentRoute *gin.RouterGroup) {
	router := parentRoute.Group("/base")
	router.Use(middleware.Auth())
	router.GET("/start_date", Base_StartDate)                                         // 7.1 本部中介-起始日期
	router.GET("/sales_achievement", Base_SalesAchievement)                           // 7.2 本部中介-销售业绩
	router.GET("/sales_profit/list", Base_SalesProfitList)                            // 7.3.1 本部中介-中介费用统计-列表
	router.GET("/sales_profit/detail", Base_SalesProfitDetail)                        // 7.3.2 本部中介-中介费用统计-详情
	router.GET("/sales_profit/setting_detail", Base_SalesProfitSettingDetail)         // 7.3.3 本部中介-中介费用统计-设置详情
	router.POST("/sales_profit/setting_modify", Base_SalesProfitSettingModify)        // 7.3.4 本部中介-中介费用统计-设置修改
	router.GET("/wait_distribution/list", Base_WaitDistributionList)                  // 7.4.1 本部中介-待分配客户-列表
	router.POST("/wait_distribution/distribution", Base_WaitDistributionDistribution) // 7.4.2 本部中介-待分配客户-分配
	router.DELETE("/wait_distribution/del/:id", Base_WaitDistributionDel)             // 7.4.3 本部中介-待分配客户-删除
	router.GET("/japan_manage/list", Base_JapanManageList)                            // 7.5.1 本部中介-日本中介管理-列表
	router.GET("/japan_manage/detail", Base_JapanManageDetail)                        // 7.5.2 本部中介-日本中介管理-详情
	router.POST("/japan_manage/add", Base_JapanManageAdd)                             // 7.5.3 本部中介-日本中介管理-添加/编辑
	router.DELETE("/japan_manage/del/:id", Base_JapanManageDel)                       // 7.5.4 本部中介-日本中介管理-删除
	router.GET("/china_manage/region_list", Base_ChinaManageRegionList)               // 7.6.1 本部中介-中国中介管理-地区列表
	router.GET("/china_manage/list", Base_ChinaManageList)                            // 7.6.2 本部中介-中国中介管理-列表
	router.GET("/china_manage/detail", Base_ChinaManageDetail)                        // 7.6.3 本部中介-中国中介管理-详情
	router.POST("/china_manage/add", Base_ChinaManageAdd)                             // 7.6.4 本部中介-中国中介管理-添加/编辑
	router.DELETE("/china_manage/del/:id", Base_ChinaManageDel)                       // 7.6.5 本部中介-中国中介管理-删除
	router.GET("/customer_manage/source_list", Base_CustomerManageSourceList)         // 7.7.1 本部中介-客户管理-来源列表
	router.GET("/customer_manage/list", Base_CustomerManageList)                      // 7.7.2 本部中介-客户管理-列表
	router.GET("/customer_manage/detail", Base_CustomerManageDetail)                  // 7.7.3 本部中介-客户管理-详情
	router.POST("/customer_manage/edit", Base_CustomerManageEdit)                     // 7.7.4 本部中介-客户管理-编辑
	router.DELETE("/customer_manage/del/:id", Base_CustomerManageDel)                 // 7.7.5 本部中介-客户管理-删除
	router.GET("/protection_period/show", Base_ProtectionPeriodShow)                  // 7.8.1 本部中介-保护期-显示
	router.POST("/protection_period/set", Base_ProtectionPeriodSet)                   // 7.8.2 本部中介-保护期-设置
	router.GET("/agency_fee/show", Base_AgencyFeeShow)                                // 7.9.1 本部中介-中介费-显示
	router.POST("/agency_fee/set", Base_AgencyFeeSet)                                 // 7.9.2 本部中介-中介费-设置
	router.POST("/notify/set", Base_NotifySet)                                        // 7.10.1 本部中介-我的通知-设置
	router.GET("/notify/list", Base_NotifyList)                                       // 7.10.2 本部中介-我的通知-列表
	router.DELETE("/notify/del/:id", Base_NotifyDel)                                  // 7.10.3 本部中介-我的通知-删除
	router.POST("/notify/marked_as_read", Base_NotifyMarkedAsRead)                    // 7.10.4 本部中介-我的通知-标记为已读
}

// 本部中介-日期列表
func Base_StartDate(c *gin.Context) {
	data, errMsg := baseModel.Base_StartDate()
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

// 本部中介-销售业绩
func Base_SalesAchievement(c *gin.Context) {
	addTime := string(c.Query("add_time"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	lastId, _ := strconv.Atoi(c.Query("last_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if addTime == "" || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允许查看销售业绩",
		})
		return
	}

	// 销售业绩列表
	data, errMsg := baseModel.Base_SalesAchievement(addTime, perPage, lastId)
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

// 本部中介-中介费用统计列表
func Base_SalesProfitList(c *gin.Context) {
	addTime := string(c.Query("add_time"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	lastId, _ := strconv.Atoi(c.Query("last_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if addTime == "" || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许查看中介费用统计",
		})
		return
	}

	// 中介费用统计列表
	data, errMsg := baseModel.Base_SalesProfitList(addTime, perPage, lastId)
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

// 本部中介-中介费用统计详情
func Base_SalesProfitDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许查看中介费用统计详情",
		})
		return
	}

	// 中介费用统计详情
	data, errMsg := baseModel.Base_SalesProfitDetail(id)
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

// 本部中介-中介费用统计设置详情
func Base_SalesProfitSettingDetail(c *gin.Context) {
	estateId, _ := strconv.Atoi(c.Query("estate_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if estateId == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许查看中介费用统计设置详情",
		})
		return
	}

	// 中介费用统计详情
	data, errMsg := baseModel.Base_SalesProfitSettingDetail(estateId)
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

// 本部中介-中介费用统计设置修改
func Base_SalesProfitSettingModify(c *gin.Context) {
	estateId, _ := strconv.Atoi(c.PostForm("estate_id"))
	agencyJson := c.PostForm("agency_json")
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if estateId == 0 || agencyJson == "" || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许修改中介费用统计设置详情",
		})
		return
	}

	// 中介费用修改
	errMsg := baseModel.Base_SalesProfitSettingModify(estateId, agencyJson)
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

// 本部中介-待分配客户列表
func Base_WaitDistributionList(c *gin.Context) {
	noticeId, _ := strconv.Atoi(c.Query("notice_id"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	lastId, _ := strconv.Atoi(c.Query("last_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许查看待分配客户列表",
		})
		return
	}

	// 待分配客户列表
	data, errMsg := baseModel.Base_WaitDistributionList(perPage, lastId)
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}

	// 标记该条消息通知为已读
	baseModel.Base_NotifyMarkedAsRead(noticeId)

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

// 本部中介-待分配客户分配
func Base_WaitDistributionDistribution(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	userId, _ := strconv.Atoi(c.PostForm("user_id"))
	leaderUserId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if id == 0 || userId == 0 || leaderUserId == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许分配客户",
		})
		return
	}

	// 分配客户
	errMsg := baseModel.Base_WaitDistributionDistribution(id, userId, leaderUserId)
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

// 本部中介-待分配客户删除
func Base_WaitDistributionDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if id == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许删除客户",
		})
		return
	}

	// 删除待分配客户
	errMsg := baseModel.Base_WaitDistributionDel(id)
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

// 本部中介-日本中介管理列表
func Base_JapanManageList(c *gin.Context) {
	keyword := c.Query("keyword")
	status, _ := strconv.Atoi(c.Query("status"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	lastId, _ := strconv.Atoi(c.Query("last_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许查看日本中介",
		})
		return
	}

	// 列表
	data, errMsg := baseModel.Base_JapanManageList(keyword, status, perPage, lastId)
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

// 本部中介-日本中介管理详情
func Base_JapanManageDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if id == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许查看日本中介详情",
		})
		return
	}

	// 详情
	data, errMsg := baseModel.Base_JapanManageDetail(id)
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

// 本部中介-日本中介管理添加/编辑
func Base_JapanManageAdd(c *gin.Context) {
	addParam := &(v1.BaseJapanManageDetailReturn{
		Id:          utils.Str2int(c.PostForm("id")),
		CompanyName: c.PostForm("company_name"),
		Address:     c.PostForm("address"),
		UserId:      utils.Str2int(c.PostForm("user_id")),
		UserName:    c.PostForm("user_name"),
		Telephone:   c.PostForm("telephone"),
		Fax:         c.PostForm("fax"),
		Email:       c.PostForm("email"),
		Password:    c.PostForm("password"),
		ExpiryDate:  c.PostForm("expiry_date"),
	})
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if addParam.CompanyName == "" || addParam.Address == "" || addParam.UserName == "" || addParam.Telephone == "" || addParam.Fax == "" || addParam.Email == "" || addParam.ExpiryDate == "" || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许添加/编辑日本中介公司",
		})
		return
	}

	// 添加/编辑
	errMsg := baseModel.Base_JapanManageAdd(addParam)
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

// 本部中介-日本中介管理删除
func Base_JapanManageDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if id == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许删除日本中介公司",
		})
		return
	}

	// 删除
	errMsg := baseModel.Base_JapanManageDel(id)
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

// 本部中介-中国中介管理地区列表
func Base_ChinaManageRegionList(c *gin.Context) {
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许查看中国中介地区",
		})
		return
	}

	// 列表
	data, errMsg := baseModel.Base_ChinaManageRegionList()
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

// 本部中介-中国中介管理列表
func Base_ChinaManageList(c *gin.Context) {
	keyword := c.Query("keyword")
	regionId, _ := strconv.Atoi(c.Query("region_id"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	lastId, _ := strconv.Atoi(c.Query("last_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许查看中国中介",
		})
		return
	}

	// 列表
	data, errMsg := baseModel.Base_ChinaManageList(keyword, regionId, perPage, lastId)
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

// 本部中介-中国中介管理详情
func Base_ChinaManageDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if id == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许查看中国中介详情",
		})
		return
	}

	// 详情
	data, errMsg := baseModel.Base_ChinaManageDetail(id)
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

// 本部中介-中国中介管理添加/编辑
func Base_ChinaManageAdd(c *gin.Context) {
	addParam := &(v1.BaseChinaManageDetailReturn{
		Id:          utils.Str2int(c.PostForm("id")),
		RegionId:    utils.Str2int(c.PostForm("region_id")),
		CompanyName: c.PostForm("company_name"),
		Address:     c.PostForm("address"),
		UserId:      utils.Str2int(c.PostForm("user_id")),
		UserName:    c.PostForm("user_name"),
		Telephone:   c.PostForm("telephone"),
		Fax:         c.PostForm("fax"),
		Email:       c.PostForm("email"),
		Password:    c.PostForm("password"),
	})
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if addParam.RegionId == 0 || addParam.CompanyName == "" || addParam.Address == "" || addParam.UserName == "" || addParam.Telephone == "" || addParam.Fax == "" || addParam.Email == "" || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许添加/编辑中国中介公司",
		})
		return
	}

	// 添加/编辑
	errMsg := baseModel.Base_ChinaManageAdd(addParam)
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

// 本部中介-中国中介管理删除
func Base_ChinaManageDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if id == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允许删除中国中介公司",
		})
		return
	}

	// 删除
	errMsg := baseModel.Base_ChinaManageDel(id)
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

// 本部中介-客户管理来源列表
func Base_CustomerManageSourceList(c *gin.Context) {
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允许查看客户管理来源列表",
		})
		return
	}

	// 列表
	data, errMsg := baseModel.Base_CustomerManageSourceList()
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

// 本部中介-客户管理列表
func Base_CustomerManageList(c *gin.Context) {
	cusParam := &(v1.BaseCustomerManageListParamater{
		Keyword:   c.Query("keyword"),
		UserId:    utils.Str2int(c.Query("user_id")),
		CompanyId: utils.Str2int(c.Query("company_id")),
		IsButt:    c.Query("is_butt"),
		IsToJapan: c.Query("is_to_japan"),
		IsAgree:   c.Query("is_agree"),
		IsPay:     c.Query("is_pay"),
		IsLoan:    c.Query("is_loan"),
		PerPage:   utils.Str2int(c.Query("per_page")),
		LastId:    utils.Str2int(c.Query("last_id")),
	})
	noticeId, _ := strconv.Atoi(c.Query("notice_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允许查看客户管理列表",
		})
		return
	}

	// 列表
	data, errMsg := baseModel.Base_CustomerManageList(cusParam)
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}

	// 标记该条消息通知为已读
	baseModel.Base_NotifyMarkedAsRead(noticeId)

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

// 本部中介-客户管理详情
func Base_CustomerManageDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if id == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允许查看客户管理详情",
		})
		return
	}

	// 列表
	data, errMsg := baseModel.Base_CustomerManageDetail(id)
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

// 本部中介-客户管理编辑
func Base_CustomerManageEdit(c *gin.Context) {
	cusParam := &(v1.BaseCustomerManageDetailReturn{
		Id:         utils.Str2int(c.PostForm("id")),
		Name:       c.PostForm("name"),
		Sex:        utils.Str2int(c.PostForm("sex")),
		Wechat:     c.PostForm("wechat"),
		IsButt:     utils.Str2int(c.PostForm("is_butt")),
		IsToJapan:  utils.Str2int(c.PostForm("is_to_japan")),
		IsAgree:    utils.Str2int(c.PostForm("is_agree")),
		IsPay:      utils.Str2int(c.PostForm("is_pay")),
		IsLoan:     utils.Str2int(c.PostForm("is_loan")),
		EstateCode: c.PostForm("estate_code"),
		Price:      c.PostForm("price"),
		UserId:     utils.Str2int(c.Request.Header.Get("user_id")),
	})
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允许编辑客户管理",
		})
		return
	}

	// 编辑
	errMsg := baseModel.Base_CustomerManageEdit(cusParam)
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

// 本部中介-客户管理删除
func Base_CustomerManageDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if id == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允许删除客户管理",
		})
		return
	}

	// 删除
	errMsg := baseModel.Base_CustomerManageDel(id)
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

// 本部中介-保护期显示
func Base_ProtectionPeriodShow(c *gin.Context) {
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允查看保护期设置",
		})
		return
	}

	// 显示
	data, errMsg := baseModel.Base_ProtectionPeriodShow()
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

// 本部中介-保护期设置
func Base_ProtectionPeriodSet(c *gin.Context) {
	protectionPeriod, _ := strconv.Atoi(c.PostForm("protection_period"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if protectionPeriod < 0 || protectionPeriod > 99 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允修改保护期设置",
		})
		return
	}

	// 设置
	errMsg := baseModel.Base_ProtectionPeriodSet(protectionPeriod)
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

// 本部中介-中介费显示
func Base_AgencyFeeShow(c *gin.Context) {
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允查看中介费设置",
		})
		return
	}

	// 显示
	data, errMsg := baseModel.Base_AgencyFeeShow()
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

// 本部中介-中介费设置
func Base_AgencyFeeSet(c *gin.Context) {
	serviceFee, _ := strconv.Atoi(c.PostForm("service_fee"))
	fixedFee, _ := strconv.Atoi(c.PostForm("fixed_fee"))
	exciseFee := c.PostForm("excise_fee")
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if serviceFee == 0 || exciseFee == "" || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && userType != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部主管不允修改中介费设置",
		})
		return
	}

	// 设置
	errMsg := baseModel.Base_AgencyFeeSet(serviceFee, fixedFee, exciseFee)
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

// 本部中介-我的通知设置
func Base_NotifySet(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允通知设置",
		})
		return
	}

	// 设置
	data, errMsg := baseModel.Base_NotifySet(userId)
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
		"data": data,
	})
	return
}

// 本部中介-我的通知列表
func Base_NotifyList(c *gin.Context) {
	status, _ := strconv.Atoi(c.Query("status"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	lastId, _ := strconv.Atoi(c.Query("last_id"))
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if userId == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允查看通知列表",
		})
		return
	}

	// 列表
	data, errMsg := baseModel.Base_NotifyList(status, perPage, lastId, userId, userType)
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
			"data": data,
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

// 本部中介-我的通知删除
func Base_NotifyDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允删除通知",
		})
		return
	}

	// 删除
	errMsg := baseModel.Base_NotifyDel(id)
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

// 本部中介-我的通知标记为已读
func Base_NotifyMarkedAsRead(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部中介不允标记为已读通知",
		})
		return
	}

	// 删除
	errMsg := baseModel.Base_NotifyMarkedAsRead(id)
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
