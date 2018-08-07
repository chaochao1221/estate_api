package v1

import (
	"estate/middleware"
	"estate/models/v1"
	"estate/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var publicModel = new(v1.PublicModel)

// 公用-路由
func Public(parentRoute *gin.RouterGroup) {
	router := parentRoute.Group("/public")
	router.GET("/japan_region_list", Public_JapanRegionList) // 3.4 公用-日本地区列表
	router.GET("/estate_detail", Public_EstateDetail)        // 3.6 公用-房源详情
	router.Use(middleware.Auth())
	router.GET("/company_detail", Public_CompanyDetail)                            // 3.1 公用-公司详情
	router.GET("/sales_manage/list", Public_SalesManageList)                       // 3.2.1 公用-销售管理-列表
	router.GET("/sales_manage/detail", Public_SalesManageDetail)                   // 3.2.2 公用-销售管理-详情
	router.POST("/sales_manage/add", Public_SalesManageAdd)                        // 3.2.3 公用-销售管理-添加/编辑
	router.DELETE("/sales_manage/del/:user_id", Public_SalesManageDel)             // 3.2.4 公用-销售管理-删除
	router.POST("/estate_manage/add", Public_EstateManageAdd)                      // 3.3.1 公用-房源管理-添加/编辑
	router.DELETE("/estate_manage/del/:estate_id", Public_EstateManageDel)         // 3.3.2 公用-房源管理-删除
	router.POST("/estate_manage/add_shelves", Public_EstateManageAddShelves)       // 3.3.3 公用-房源管理-上架
	router.POST("/estate_manage/remove_shelves", Public_EstateManageRemoveShelves) // 3.3.4 公用-房源管理-下架
	router.GET("/estate_list", Public_EstateList)                                  // 3.5 公用-房源列表
	router.POST("/feedback", Public_Feedback)                                      // 3.7 公用-意见反馈
	router.POST("/contact", Public_Contact)                                        // 3.8 公用-联系方式
}

// 公用-公司详情
func Public_CompanyDetail(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
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
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	if userId == 0 {
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
	data, errMsg := publicModel.Public_SalesManageList(userId)
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

// 公用-销售管理添加/编辑
func Public_SalesManageAdd(c *gin.Context) {
	leaderUserId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	userId, _ := strconv.Atoi(c.PostForm("user_id")) // 添加非必传，编辑必传
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password") // 添加必传，编辑非必传
	if leaderUserId == 0 || name == "" || email == "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if userType == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "该账户非公司主管，不允许添加/编辑员工",
		})
		return
	}

	// 添加/编辑
	errMsg := publicModel.Public_SalesManageAdd(leaderUserId, userId, name, email, password)
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

// 公用-销售管理删除
func Public_SalesManageDel(c *gin.Context) {
	leaderUserId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	userId, _ := strconv.Atoi(c.Params.ByName("user_id")) // 添加非必传，编辑必传
	if groupId == 0 || userId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if userType == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "该账户非公司主管，不允许删除员工",
		})
		return
	}

	// 删除
	errMsg := publicModel.Public_SalesManageDel(leaderUserId, groupId, userId)
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

// 公用-房源管理添加/编辑
func Public_EstateManageAdd(c *gin.Context) {
	estparam := &(v1.PublicEstateManageAddParameter{
		EstateId:          utils.Str2int(c.PostForm("estate_id")),
		Price:             utils.Str2int(c.PostForm("price")),
		Points:            utils.Str2int(c.PostForm("points")),
		Huxing:            c.PostForm("huxing"),
		MeasureArea:       c.PostForm("measure_area"),
		HousingType:       utils.Str2int(c.PostForm("housing_type")),
		Floor:             utils.Str2int(c.PostForm("floor")),
		TotalFloor:        utils.Str2int(c.PostForm("total_floor")),
		BuildingTime:      c.PostForm("building_time"),
		BuildingStructure: utils.Str2int(c.PostForm("building_structure")),
		LandRights:        utils.Str2int(c.PostForm("land_rights")),
		Orientation:       c.PostForm("orientation"),
		State:             utils.Str2int(c.PostForm("state")),
		Rent:              utils.Str2int(c.PostForm("rent")),
		ReturnRate:        c.PostForm("return_rate"),
		RepairFee:         utils.Str2int(c.PostForm("repair_fee")),
		ManageFee:         utils.Str2int(c.PostForm("manage_fee")),
		RegionId:          utils.Str2int(c.PostForm("region_id")),
		Traffic:           c.PostForm("traffic"),
		Address:           c.PostForm("address"),
		Picture:           c.PostForm("picture"),
		UserId:            utils.Str2int(c.Request.Header.Get("user_id")),
	})
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if estparam.UserId == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}
	if groupId != 1 && groupId != 3 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "非本部、日本中介，不允许发布房源",
		})
		return
	}

	// 添加/编辑
	errMsg := publicModel.Public_EstateManageAdd(estparam)
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

// 公用-房源管理删除
func Public_EstateManageDel(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	estateId, _ := strconv.Atoi(c.Params.ByName("estate_id"))
	if userId == 0 || groupId == 0 || estateId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 判断是否存在删除权限
	_, errMsg := publicModel.ExistEstateManagePermissions(&(v1.EstateManagePermissionsParamater{
		GroupId:  groupId,
		UserId:   userId,
		UserType: userType,
		EstateId: estateId,
	}))
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}

	// 删除
	errMsg = publicModel.Public_EstateManageDel(estateId)
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

// 公用-房源管理上架
func Public_EstateManageAddShelves(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	estateId, _ := strconv.Atoi(c.PostForm("estate_id"))
	if userId == 0 || groupId == 0 || estateId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 判断是否存在删除权限
	_, errMsg := publicModel.ExistEstateManagePermissions(&(v1.EstateManagePermissionsParamater{
		GroupId:  groupId,
		UserId:   userId,
		UserType: userType,
		EstateId: estateId,
	}))
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}

	// 上架
	errMsg = publicModel.Public_EstateManageAddShelves(estateId)
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

// 公用-房源管理下架
func Public_EstateManageRemoveShelves(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	userType, _ := strconv.Atoi(c.Request.Header.Get("user_type"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	estateId, _ := strconv.Atoi(c.PostForm("estate_id"))
	if userId == 0 || groupId == 0 || estateId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 判断是否存在删除权限
	_, errMsg := publicModel.ExistEstateManagePermissions(&(v1.EstateManagePermissionsParamater{
		GroupId:  groupId,
		UserId:   userId,
		UserType: userType,
		EstateId: estateId,
	}))
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}

	// 下架
	errMsg = publicModel.Public_EstateManageRemoveShelves(estateId)
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

// 公用-日本地区列表
func Public_JapanRegionList(c *gin.Context) {
	// 日本地区列表
	data, errMsg := publicModel.Public_JapanRegionList()
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

// 公用-房源列表（非游客使用）
func Public_EstateList(c *gin.Context) {
	estParam := &(v1.PublicEstateListParamter{
		Keyword:    c.Query("keyword"),
		Listorder:  utils.Str2int(c.Query("listorder")),
		ScreenJson: c.Query("screen_json"),
		Status:     utils.Str2int(c.Query("status")),
		PerPage:    utils.Str2int(c.Query("per_page")),
		LastId:     utils.Str2int(c.Query("last_id")),
		UserId:     utils.Str2int(c.Request.Header.Get("user_id")),
		UserType:   utils.Str2int(c.Request.Header.Get("user_type")),
		GroupId:    utils.Str2int(c.Request.Header.Get("group_id")),
	})
	if estParam.UserId == 0 || estParam.GroupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 房源列表
	data, errMsg := publicModel.Public_EstateList(estParam)
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

// 公用-房源详情
func Public_EstateDetail(c *gin.Context) {
	estateId, _ := strconv.Atoi(c.Query("estate_id"))
	if estateId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 房源详情
	data, errMsg := publicModel.Public_EstateDetail(estateId)
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

// 公用-意见反馈
func Public_Feedback(c *gin.Context) {
	types, _ := strconv.Atoi(c.PostForm("type"))
	contact := c.PostForm("contact")
	content := c.PostForm("content")
	if contact == "" || content == "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 意见反馈
	errMsg := publicModel.Public_Feedback(types, contact, content)
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

// 公用-联系方式
func Public_Contact(c *gin.Context) {
	estateId, _ := strconv.Atoi(c.Query("estate_id"))
	if estateId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 联系方式
	data, errMsg := publicModel.Public_Contact(estateId)
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
