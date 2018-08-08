package v1

import (
	"encoding/json"
	"estate/db"
	"estate/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golibs/uuid"
)

type PublicModel struct{}

var (
	userModel = new(UserModel)
	baseModel = new(BaseModel)
)

type PublicCompanyDetailReturn struct {
	GroupId         int    `json:"group_id"`
	UserId          int    `json:"user_id"`
	UserType        int    `json:"user_type"`
	CompanyName     string `json:"company_name"`
	RecommendNumber int    `json:"recommend_number"`
	ReleaseNumber   int    `json:"release_number"`
	ButtNumber      int    `json:"butt_number"`
	DealNumber      int    `json:"deal_number"`
	ExpiryDate      string `json:"expiry_date"`
}

// 公用-公司详情
func (this *PublicModel) Public_CompanyDetail(userId int) (data *PublicCompanyDetailReturn, errMsg string) {
	// 用户信息
	userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: userId})
	if errMsg != "" {
		return data, errMsg
	}

	// 公司信息
	companyInfo, errMsg := userModel.GetCompanyInfo(userInfo.CompanyId)
	if errMsg != "" {
		return data, errMsg
	}

	return &PublicCompanyDetailReturn{
		GroupId:         companyInfo.GroupId,
		UserId:          userId,
		UserType:        userInfo.UserType,
		CompanyName:     companyInfo.Name,
		RecommendNumber: companyInfo.RecommendNumber,
		ReleaseNumber:   companyInfo.ReleaseNumber,
		ButtNumber:      companyInfo.ButtNumber,
		DealNumber:      companyInfo.DealNumber,
		ExpiryDate:      companyInfo.ExpiryDate,
	}, ""
}

type PublicSalesManageListReturn struct {
	List []SalesManageList `json:"list"`
}

type SalesManageList struct {
	UserId  int    `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	AddTime string `json:"add_time"`
	Mark    int    `json:"mark"`
}

// 公用-销售管理列表
func (this *PublicModel) Public_SalesManageList(userId int) (data *PublicSalesManageListReturn, errMsg string) {
	data = new(PublicSalesManageListReturn)

	// 用户信息
	userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: userId})
	if errMsg != "" {
		return data, errMsg
	}

	// 获取已经分配的销售
	sql := `SELECT DISTINCT user_id FROM base_distribution`
	rows, err := db.Db.Query(sql)
	if err != nil {
		return data, "获取已经分配的销售失败"
	}
	userIdMap := make(map[int]int)
	if len(rows) > 0 {
		for _, value := range rows {
			userIdMap[utils.Str2int(string(value["user_id"]))] = 1
		}
	}

	// 获取公司员工列表
	sql = `SELECT id, name, email, add_time
		   FROM p_user
		   WHERE company_id=? AND user_type=0`
	rows, err = db.Db.Query(sql, userInfo.CompanyId)
	if err != nil {
		return data, "获取员工列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for _, value := range rows {
		data.List = append(data.List, SalesManageList{
			UserId:  utils.Str2int(string(value["id"])),
			Name:    string(value["name"]),
			Email:   string(value["email"]),
			AddTime: DateFormat(string(value["add_time"]), "2006/01/02"),
			Mark:    userIdMap[utils.Str2int(string(value["id"]))],
		})
	}
	return
}

/*
* @Title DateFormat
* @Description 日期格式化
* @Parameter add_time string
* @Parameter format string
* @Return data string
 */
func DateFormat(add_time, format string) (data string) {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", add_time, time.Local)
	return t.Format(format)
}

type PublicSalesManageDetailReturn struct {
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// 公用-销售管理详情
func (this *PublicModel) Public_SalesManageDetail(userId int) (data *PublicSalesManageDetailReturn, errMsg string) {
	// 用户信息
	userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: userId})
	if errMsg != "" {
		return data, errMsg
	}
	if userInfo == nil {
		return data, "该用户不存咋"
	}

	return &PublicSalesManageDetailReturn{
		UserId: userId,
		Name:   userInfo.Name,
		Email:  userInfo.Email,
	}, ""
}

// 公用-销售管理添加/编辑
func (this *PublicModel) Public_SalesManageAdd(leaderUserId, userId int, name, email, password string) (errMsg string) {
	// 判断除自己以外是否还存在该邮箱
	_, errMsg = this.ExistEmail(email, userId)
	if errMsg != "" {
		return
	}

	// 用户信息
	userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: leaderUserId})
	if errMsg != "" {
		return errMsg
	}

	// 根据userId是否为0来判断是添加还是编辑
	if userId == 0 { // 添加
		// 添加
		sql := `INSERT INTO p_user(company_id, email, password, name) VALUES(?,?,?,?)`
		_, err := db.Db.Exec(sql, userInfo.CompanyId, email, string(utils.HashPassword(password)), name)
		if err != nil {
			return "添加用户失败"
		}
	} else { // 编辑
		var passwordSql string
		if password != "" {
			passwordSql = ` ,password="` + string(utils.HashPassword(password)) + `"`
		}
		sql := `UPDATE p_user
			   SET email=?, name=? ` + passwordSql +
			`WHERE id=?`
		_, err := db.Db.Exec(sql, email, name, userId)
		if err != nil {
			return "编辑用户失败"
		}
	}

	return
}

/*
* @Title ExistEmail
* @Description 判断除自己以外是否还存在该邮箱
* @Parameter email string
* @Parameter userId int
* @Return data bool
* @Return errMsg string
 */
func (this *PublicModel) ExistEmail(email string, userId int) (data bool, errMsg string) {
	fmt.Println(email, userId)
	sql := `SELECT id FROM p_user WHERE email=? AND id<>?`
	row, err := db.Db.Query(sql, email, userId)
	if err != nil {
		return false, "获取邮箱失败"
	}
	if len(row) > 0 {
		return false, "The email has already existed"
	}
	return true, ""
}

// 公用-销售管理删除
func (this *PublicModel) Public_SalesManageDel(leaderUserId, groupId, userId int) (errMsg string) {
	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 中介删除（本部、中国、日本）
	// 1）中介用户身份物理删除
	// 1）中国中介、日本中介相关数据挂到该公司主管上
	// 2）本部中介房源信息挂到本部主管上，本部中介客户信息进入待分配状态。销售业绩、中介费用统计中没有中介关系的数据归入其他
	sql := `DELETE FROM p_user WHERE id=?`
	_, err := transaction.Exec(sql, userId)
	if err != nil {
		transaction.Rollback()
		return "删除用户身份失败"
	}

	switch groupId {
	case 1: // 本部
		// 更改房源发布者
		sql = `UPDATE p_estate SET user_id=? WHERE user_id=?`
		_, err = transaction.Exec(sql, leaderUserId, userId)
		if err != nil {
			transaction.Rollback()
			return "更新本部发布者失败"
		}

		// 更新是否分配状态
		sql = `UPDATE p_recommend r
			   LEFT JOIN base_distribution d ON d.recommend_id=r.id
			   SET r.is_distribution=0
			   WHERE d.user_id=?`
		_, err = transaction.Exec(sql, userId)
		if err != nil {
			transaction.Rollback()
			return "更新分配状态失败"
		}

		// 删除分配关系
		sql = `DELETE FROM base_distribution WHERE user_id=?`
		_, err = transaction.Exec(sql, userId)
		if err != nil {
			transaction.Rollback()
			return "删除分配状态失败"
		}
	case 2: // 中国
		// 更改推荐者
		sql = `UPDATE p_recommend SET user_id=? WHERE user_id=?`
		_, err = transaction.Exec(sql, leaderUserId, userId)
		if err != nil {
			transaction.Rollback()
			return "更新推荐者失败"
		}
	case 3: // 日本
		// 更改房源发布者
		sql = `UPDATE p_estate SET user_id=? WHERE user_id=?`
		_, err = transaction.Exec(sql, leaderUserId, userId)
		if err != nil {
			transaction.Rollback()
			return "更新日本中介发布角色失败"
		}
	default:
		return "不存在该中介分组"
	}

	// 提交事务
	if err := transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}
	return
}

type EstateManagePermissionsParamater struct {
	GroupId  int
	UserId   int
	UserType int
	EstateId int
}

/*
* @Title ExistEstateManagePermissions
* @Description 是否存在房源管理权限
* @Parameter e *GetUserInfoParameter
* @Return data bool
* @Return errMsg string
 */
func (this *PublicModel) ExistEstateManagePermissions(e *EstateManagePermissionsParamater) (data bool, errMsg string) {
	// 判断该用户是管理员还是普通销售，管理员可以删除本中介公司下的所有房源，销售只能删除自己发布的房源，已成交的房源不允许操作
	sql := `SELECT e.user_id, u.company_id, e.status
			FROM p_estate e
			LEFT JOIN p_user u ON u.id=e.user_id
			WHERE e.id=? AND e.is_del=0`
	row, err := db.Db.Query(sql, e.EstateId)
	if err != nil {
		return false, "获取房源发布者失败"
	}
	if len(row) == 0 {
		return false, "该房源不存在"
	}
	if utils.Str2int(string(row[0]["status"])) == 3 {
		return false, "该房源已成交，不允许操作"
	}

	switch e.UserType {
	case 0: // 销售
		if e.UserId != utils.Str2int(string(row[0]["user_id"])) {
			return false, "该用户不是此房源的发布者，不允许操作该房源"
		}
	case 1: // 管理员
		// 用户信息
		userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: e.UserId})
		if errMsg != "" {
			return false, errMsg
		}
		if userInfo.CompanyId != utils.Str2int(string(row[0]["company_id"])) {
			return false, "该用户不是此房源公司的管理员，不允许操作该房源"
		}
	default:
		return false, "不存在该用户类型"
	}
	return true, ""
}

type PublicEstateManageAddParameter struct {
	EstateId          int
	Price             int
	Points            int
	Huxing            string
	MeasureArea       string
	HousingType       int
	Floor             int
	TotalFloor        int
	BuildingTime      string
	BuildingStructure int
	LandRights        int
	Orientation       string
	State             int
	Rent              int
	ReturnRate        string
	RepairFee         int
	ManageFee         int
	RegionId          int
	Traffic           string
	Address           string
	Picture           string
	UserId            int
}

// 公用-房源管理添加/编辑
func (this *PublicModel) Public_EstateManageAdd(estParam *PublicEstateManageAddParameter) (errMsg string) {
	// 中介费
	agencyFee, errMsg := this.GetAgencyFee(estParam.EstateId, estParam.Price)
	if errMsg != "" {
		return errMsg
	}

	// 用户信息
	userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: estParam.UserId})
	if errMsg != "" {
		return errMsg
	}

	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 更新房产
	if estParam.EstateId == 0 { // 新增
		// 新增房源
		sql := `INSERT INTO p_estate(user_id, code, agency_fee, price, points, measure_area, huxing, housing_type, total_floor, floor, building_time, building_structure, land_rights, orientation, state, rent, return_rate, repair_fee, manage_fee, region_id, traffic, address, picture) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		_, err := transaction.Exec(sql, estParam.UserId, uuid.Rand().Hex(), agencyFee, estParam.Price, estParam.Points, estParam.MeasureArea, estParam.Huxing, estParam.HousingType, estParam.TotalFloor, estParam.Floor, estParam.BuildingTime, estParam.BuildingStructure, estParam.LandRights, estParam.Orientation, estParam.State, estParam.Rent, estParam.ReturnRate, estParam.RepairFee, estParam.ManageFee, estParam.RegionId, estParam.Traffic, estParam.Address, estParam.Picture)
		if err != nil {
			transaction.Rollback()
			return "新增房产失败"
		}

		// 更新发布房源数量
		sql = `UPDATE p_company
			   SET release_number=release_number+1
			   WHERE id=?`
		_, err = transaction.Exec(sql, userInfo.CompanyId)
		if err != nil {
			transaction.Rollback()
			return "更新发布房源数量失败"
		}
	} else { // 编辑
		// 编辑房源
		sql := `UPDATE p_estate
				SET agency_fee=?, price=?, points=?, measure_area=?, huxing=?, housing_type=?, total_floor=?, floor=?, building_time=?, building_structure=?, land_rights=?, orientation=?, state=?, rent=?, return_rate=?, repair_fee=?, manage_fee=?, region_id=?, traffic=?, address=?, picture=?
				WHERE id=?`
		_, err := transaction.Exec(sql, agencyFee, estParam.Price, estParam.Points, estParam.MeasureArea, estParam.Huxing, estParam.HousingType, estParam.TotalFloor, estParam.Floor, estParam.BuildingTime, estParam.BuildingStructure, estParam.LandRights, estParam.Orientation, estParam.State, estParam.Rent, estParam.ReturnRate, estParam.RepairFee, estParam.ManageFee, estParam.RegionId, estParam.Traffic, estParam.Address, estParam.Picture, estParam.EstateId)
		if err != nil {
			transaction.Rollback()
			return "更新房产失败"
		}
	}

	// 提交事务
	if err := transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}

	return
}

/*
* @Title GetAgencyFee
* @Description 中介费
* @Parameter price int
* @Return data int
* @Return errMsg string
 */
func (this *PublicModel) GetAgencyFee(estateId, price int) (data int, errMsg string) {
	if estateId == 0 {
		// 本部基础信息
		baseInfo, errMsg := userModel.GetBaseInfo()
		if errMsg != "" {
			return data, errMsg
		}
		data = utils.Float64Toint(float64(price)*(float64(baseInfo.ServiceFee)/100)*(utils.Str2float64(baseInfo.ExciseFee)/100) + float64(baseInfo.FixedFee))
	} else {
		// 中介费详情
		profit, errMsg := baseModel.Base_SalesProfitSettingDetail(estateId)
		if errMsg != "" {
			return data, errMsg
		}
		data = utils.Float64Toint(float64(price)*(float64(profit.Buyer.ServiceFee)/100)*(utils.Str2float64(profit.Buyer.ExciseFee)/100)+float64(profit.Buyer.FixedFee)) +
			utils.Float64Toint(float64(price)*(float64(profit.Seller.ServiceFee)/100)*(utils.Str2float64(profit.Seller.ExciseFee)/100)+float64(profit.Seller.FixedFee))
	}

	return
}

// 公用-房源管理删除
func (this *PublicModel) Public_EstateManageDel(estateId int) (errMsg string) {
	sql := `UPDATE p_estate SET is_del=1 WHERE id=?`
	_, err := db.Db.Exec(sql, estateId)
	if err != nil {
		return "删除房源失败"
	}
	return
}

// 公用-房源管理上架
func (this *PublicModel) Public_EstateManageAddShelves(estateId int) (errMsg string) {
	// 判断该房源是否处于下架状态
	sql := `SELECT status FROM p_estate WHERE id=?`
	row, err := db.Db.Query(sql, estateId)
	if err != nil {
		return "获取房源状态失败"
	}
	if utils.Str2int(string(row[0]["status"])) != 2 {
		return "该房源不处于下架状态，不允许上架"
	}

	// 上架
	sql = `UPDATE p_estate SET status=1 WHERE id=?`
	_, err = db.Db.Exec(sql, estateId)
	if err != nil {
		return "房源上架失败"
	}
	return
}

// 公用-房源管理下架
func (this *PublicModel) Public_EstateManageRemoveShelves(estateId int) (errMsg string) {
	// 判断该房源是否处于上架状态
	sql := `SELECT status FROM p_estate WHERE id=?`
	row, err := db.Db.Query(sql, estateId)
	if err != nil {
		return "获取房源状态失败"
	}
	if utils.Str2int(string(row[0]["status"])) != 1 {
		return "该房源不处于上架状态，不允许下架"
	}

	// 上架
	sql = `UPDATE p_estate SET status=2 WHERE id=?`
	_, err = db.Db.Exec(sql, estateId)
	if err != nil {
		return "房源上架失败"
	}
	return
}

type PublicJapanRegionListReturn struct {
	Region []JapanRegionList `json:"region"`
}

type JapanRegionList struct {
	RegionId   int             `json:"region_id"`
	RegionName string          `json:"region_name"`
	Area       []JapanAreaList `json:"area"`
}

type JapanAreaList struct {
	AreaId   int    `json:"area_id"`
	AreaName string `json:"area_name"`
}

// 公用-日本地区列表
func (this *PublicModel) Public_JapanRegionList() (data *PublicJapanRegionListReturn, errMsg string) {
	data = new(PublicJapanRegionListReturn)

	// 地区列表
	sql := `SELECT id, p_id, name, type
			FROM p_region
			WHERE group_id<>2`
	rows, err := db.Db.Query(sql)
	if err != nil {
		return data, "获取地区列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	var (
		regionArray   = make([]JapanRegionList, 0)
		regionMap     = make(map[int]int)
		regionAreaMap = make(map[int][]JapanAreaList)
	)
	for _, value := range rows {
		id, _ := strconv.Atoi(string(value["id"]))
		typeId, _ := strconv.Atoi(string(value["type"]))
		pId, _ := strconv.Atoi(string(value["p_id"]))
		if typeId == 1 { // 市
			if _, ok := regionMap[id]; !ok {
				regionArray = append(regionArray, JapanRegionList{
					RegionId:   id,
					RegionName: string(value["name"]),
				})
			}
		} else { // 区
			regionAreaMap[pId] = append(regionAreaMap[pId], JapanAreaList{
				AreaId:   id,
				AreaName: string(value["name"]),
			})
		}
	}
	for _, value := range regionArray {
		value.Area = regionAreaMap[value.RegionId]
		if len(value.Area) == 0 {
			value.Area = []JapanAreaList{}
		}
		data.Region = append(data.Region, value)
	}
	return
}

type PublicEstateListParamter struct {
	Keyword    string
	Listorder  int
	ScreenJson string
	Status     int
	PerPage    int
	LastId     int
	UserId     int
	UserType   int
	GroupId    int
}

type ScreenJson struct {
	RegionId          int
	AreaId            int
	PriceMin          int
	PriceMax          int
	AreaMin           string
	AreaMax           string
	RoomNumber        int
	HousingType       string
	BuildingTime      int
	Orientation       string
	Floor             int
	BuildingStructure int
}

type PublicEstateListReturn struct {
	List       []PublicEstateList `json:"list"`
	Pagenation Pagenation         `json:"pagenation"`
}

type PublicEstateList struct {
	EstateId    int    `json:"estate_id"`
	EstateName  string `json:"estate_name"`
	Huxing      string `json:"huxing"`
	HuxingAlias string `json:"huxing_alias"`
	MeasureArea string `json:"measure_area"`
	RegionName  string `json:"region_name"`
	HousingType int    `json:"housing_type"`
	Price       string `json:"price"`
	PriceRmb    string `json:"price_rmb"`
}

// 公用-房源列表
func (this *PublicModel) Public_EstateList(estParam *PublicEstateListParamter) (data *PublicEstateListReturn, errMsg string) {
	data = new(PublicEstateListReturn)
	data.Pagenation.LastId = -1

	// 条件
	var where, order string

	// 关键字
	if estParam.Keyword != "" {
		where = ` AND e.code="` + estParam.Keyword + `" `
	}

	// 筛选
	if estParam.ScreenJson != "" {
		screen := new(ScreenJson)
		if err := json.Unmarshal([]byte(estParam.ScreenJson), screen); err != nil {
			return data, "解析筛选json失败"
		}

		// 地区
		if screen.AreaId > 0 {
			where += ` e.region_id=` + strconv.Itoa(screen.AreaId)
		} else if screen.RegionId > 0 {
			sql := `SELECT id FROM p_region WHERE p_id=?`
			rows, err := db.Db.Query(sql, screen.RegionId)
			if err != nil {
				return data, "获取地区失败"
			}
			areaIdArray := make([]string, 0)
			for _, value := range rows {
				areaIdArray = append(areaIdArray, string(value["id"]))
			}
			where += ` e.region_id IN(` + strings.Join(areaIdArray, ",") + `)`
		}

		// 价格区间
		if screen.PriceMin > 0 {
			where += ` e.price>=` + strconv.Itoa(screen.PriceMin)
		}
		if screen.PriceMax > 0 {
			where += ` e.price<=` + strconv.Itoa(screen.PriceMax)
		}

		// 面积区间
		if screen.AreaMin != "" {
			where += ` e.measure_area>=` + screen.AreaMin
		}
		if screen.AreaMax != "" {
			where += ` e.measure_area<=` + screen.AreaMax
		}

		// 房型
		if screen.RoomNumber > 0 {
			where += ` e.huxing REGEXP '^` + strconv.Itoa(screen.RoomNumber) + `'`
		}

		// 类型
		switch screen.HousingType {
		case "0":
			where += ` e.housing_type NOT IN(1,2,3,4)`
		case "1", "2", "3", "4":
			where += ` e.housing_type IN(1,2,3,4)`
		}

		// 房龄
		switch screen.BuildingTime {
		case 1: // 2年以下
			where += ` ADN date_format(date_sub(curdate(), interval 2 year),'%Y-%m')<e.building_time`
		case 2: // 2-5年
			where += ` ADN date_format(date_sub(curdate(), interval 5 year),'%Y-%m')<=e.building_time AND date_format(date_sub(curdate(), interval 2 year),'%Y-%m')>=e.building_time`
		case 3: // 5-10年
			where += ` ADN date_format(date_sub(curdate(), interval 10 year),'%Y-%m')<=e.building_time AND date_format(date_sub(curdate(), interval 5 year),'%Y-%m')>=e.building_time`
		case 4: // 10年以上
			where += ` ADN date_format(date_sub(curdate(), interval 10 year),'%Y-%m')>e.building_time`
		}

		// 朝向
		if screen.Orientation != "" {
			where += ` AND e.orientation="` + screen.Orientation + `"`
		}

		// 楼层
		switch screen.Floor {
		case 1:
			where += ` AND e.floor<6`
		case 2:
			where += ` AND e.floor<=12 AND e.floor>=6`
		case 3:
			where += ` AND e.floor>12`
		}

		// 建筑结构
		if screen.BuildingStructure > 0 {
			where += ` AND e.building_structure=` + strconv.Itoa(screen.BuildingStructure)
		}
	}

	// 状态
	if estParam.UserType == 1 { // 主管
		// 用户信息
		userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: estParam.UserId})
		if errMsg != "" {
			return data, "获取用户信息失败"
		}

		switch estParam.Status {
		case 1: // 我的
			where += ` AND u.company_id=` + strconv.Itoa(userInfo.CompanyId)
		case 2: // 已下架
			where += ` AND u.company_id=` + strconv.Itoa(userInfo.CompanyId) + ` AND e.status=2`
		case 3: // 已成交
			where += ` AND u.company_id=` + strconv.Itoa(userInfo.CompanyId) + ` AND e.status=3`
		}
	} else { // 销售
		switch estParam.Status {
		case 1: // 我的
			where += ` AND e.user_id=` + strconv.Itoa(estParam.UserId)
		case 2: // 已下架
			where += ` AND e.user_id=` + strconv.Itoa(estParam.UserId) + ` AND e.status=2`
		case 3: // 已成交
			where += ` AND e.user_id=` + strconv.Itoa(estParam.UserId) + ` AND e.status=3`
		}
	}

	// 分页
	if estParam.PerPage == 0 {
		estParam.PerPage = 10
	}

	// 排序
	switch estParam.Listorder {
	case 0: // 最新
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT add_time FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源发布时间失败"
			}
			addTime := string(row[0]["add_time"])
			where += ` AND (e.add_time<"` + addTime + `" OR (e.add_time="` + addTime + `" AND e.id<` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.add_time DESC, e.id DESC `
	case 1: // 面积正序
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT measure_area FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源面积失败"
			}
			measureArea := string(row[0]["measure_area"])
			where += ` AND (e.measure_area>"` + measureArea + `" OR (e.measure_area="` + measureArea + `" AND e.id>` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.measure_area, e.id`
	case 2: // 面积倒序
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT measure_area FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源面积失败"
			}
			measureArea := string(row[0]["measure_area"])
			where += ` AND (e.measure_area<"` + measureArea + `" OR (e.measure_area="` + measureArea + `" AND e.id<` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.measure_area DESC, e.id DESC`
	case 3: // 价格正序
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT price FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源面积失败"
			}
			price := string(row[0]["price"])
			where += ` AND (e.price>"` + price + `" OR (e.price="` + price + `" AND e.id>` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.price, e.id`
	case 4: // 价格倒序
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT price FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源面积失败"
			}
			price := string(row[0]["price"])
			where += ` AND (e.price<"` + price + `" OR (e.price="` + price + `" AND e.id<` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.price DESC, e.id DESC`
	default:
		return data, "排序类型错误"
	}

	// 列表
	sql := `SELECT e.id, e.huxing, e.measure_area, e.price, e.housing_type, r.p_id regionPId, r.name regionName, r.type regionType
			FROM p_estate e
			LEFT JOIN p_region r ON r.id=e.region_id
			LEFT JOIN p_user u ON u.id=e.user_id
			WHERE e.is_del=0 ` + where + order + ` LIMIT 0,?`
	rows, err := db.Db.Query(sql, estParam.PerPage+1)
	fmt.Println(sql)
	if err != nil {
		return data, "获取房源列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for key, value := range rows {
		if key < estParam.PerPage {
			// 地区名称
			regionName, areaName, errMsg := this.GetEstateName(&GetRegionAreaNameParameter{
				RegionPId:  utils.Str2int(string(value["regionPId"])),
				RegionType: utils.Str2int(string(value["regionType"])),
				RegionName: string(value["regionName"]),
			})
			if errMsg != "" {
				return data, errMsg
			}

			data.List = append(data.List, PublicEstateList{
				EstateId:    utils.Str2int(string(value["id"])),
				EstateName:  regionName + areaName + GetHousingTypeName(utils.Str2int(string(value["housing_type"]))),
				Huxing:      string(value["huxing"]),
				HuxingAlias: GetHuxingAlias(string(value["huxing"])),
				MeasureArea: string(value["measure_area"]),
				RegionName:  regionName,
				HousingType: utils.Str2int(string(value["housing_type"])),
				Price:       string(value["price"]),
				PriceRmb:    "",
			})

			estParam.LastId = utils.Str2int(string(value["id"]))
		} else {
			data.Pagenation.LastId = estParam.LastId
		}
	}

	return
}

type PublicEstateDetailReturn struct {
	EstateId          int    `json:"estate_id"`
	EstateName        string `json:"estate_name"`
	EstateCode        string `json:"estate_code"`
	Price             string `json:"price"`
	PriceRmb          string `json:"price_rmb"`
	Huxing            string `json:"huxing"`
	HuxingAlias       string `json:"huxing_alias"`
	MeasureArea       string `json:"measure_area"`
	HousingType       int    `json:"housing_type"`
	LandRights        int    `json:"land_rights"`
	BuildingTime      string `json:"building_time"`
	Floor             int    `json:"floor"`
	TotalFloor        int    `json:"total_floor"`
	BuildingStructure int    `json:"building_structure"`
	Orientation       string `json:"orientation"`
	RepairFee         int    `json:"repair_fee"`
	State             int    `json:"state"`
	Rent              int    `json:"rent"`
	ReturnRate        string `json:"return_rate"`
	ManageFee         int    `json:"manage_fee"`
	RegionName        string `json:"region_name"`
}

// 公用-房源详情
func (this *PublicModel) Public_EstateDetail(estateId int) (data *PublicEstateDetailReturn, errMsg string) {
	sql := `SELECT e.id, e.code, e.room_number, e.exist_living_room, e.exist_dining_room,
			e.exist_kitchen, e.measure_area, e.housing_type, e.price, e.land_rights, e.building_time, e.floor, e.total_floor,
			e.building_structure, e.orientation, e.repair_fee, e.state, e.rent, e.return_rate, e.manage_fee, r.type region_type,
			FROM p_estate e
			LEFT JOIN p_region r ON r.id=e.region_id
			WHERE e.id=? AND e.status=1 AND e.is_del=0 AND e.add_time<?`
	row, err := db.Db.Query(sql, estateId)
	if err != nil {
		return data, "获取房源信息失败"
	}
	if len(row) == 0 {
		return data, "该房源信息不存在"
	}

	// 地区名称
	regionName, areaName, errMsg := this.GetEstateName(&GetRegionAreaNameParameter{
		RegionPId:  utils.Str2int(string(row[0]["regionPId"])),
		RegionType: utils.Str2int(string(row[0]["regionType"])),
		RegionName: string(row[0]["regionName"]),
	})
	if errMsg != "" {
		return data, errMsg
	}

	// 建筑年月
	buildingTimeArray := strings.Split(string(row[0]["building_time"]), "-")
	buildingTime := buildingTimeArray[0] + "年" + utils.Int2str(utils.Str2int(buildingTimeArray[1])) + "月"

	// 封装数据
	return &PublicEstateDetailReturn{
		EstateId:          utils.Str2int(string(row[0]["id"])),
		EstateName:        regionName + areaName + GetHousingTypeName(utils.Str2int(string(row[0]["housing_type"]))),
		EstateCode:        string(row[0]["code"]),
		Price:             string(row[0]["price"]),
		PriceRmb:          "",
		Huxing:            string(row[0]["huxing"]),
		HuxingAlias:       GetHuxingAlias(string(row[0]["huxing"])),
		MeasureArea:       string(row[0]["measure_area"]),
		HousingType:       utils.Str2int(string(row[0]["housing_type"])),
		LandRights:        utils.Str2int(string(row[0]["land_rights"])),
		BuildingTime:      buildingTime,
		Floor:             utils.Str2int(string(row[0]["floor"])),
		TotalFloor:        utils.Str2int(string(row[0]["total_floor"])),
		BuildingStructure: utils.Str2int(string(row[0]["building_structure"])),
		Orientation:       string(row[0]["orientation"]),
		RepairFee:         utils.Str2int(string(row[0]["repair_fee"])),
		State:             utils.Str2int(string(row[0]["state"])),
		Rent:              utils.Str2int(string(row[0]["rent"])),
		ReturnRate:        string(row[0]["return_rate"]),
		ManageFee:         utils.Str2int(string(row[0]["manage_fee"])),
		RegionName:        regionName + `-` + areaName,
	}, ""
}

type GetRegionAreaNameParameter struct {
	RegionPId  int
	RegionType int
	RegionName string
}

/*
* @Title GetRegionAreaName
* @Description 地区名称
* @Parameter estParam *GetEstateNameParameter
* @Return data string
* @Return errMsg string
 */
func (this *PublicModel) GetEstateName(estParam *GetRegionAreaNameParameter) (regionName, areaName, errMsg string) {
	if estParam.RegionType == 1 { // 市
		regionName = estParam.RegionName
	} else { // 区
		sql := `SELECT name FROM p_region WHERE id=?`
		row, err := db.Db.Query(sql, estParam.RegionPId)
		if err != nil {
			return regionName, areaName, "获取地区失败"
		}
		regionName, areaName = string(row[0]["name"]), estParam.RegionName
	}
	return
}

/*
* @Title GetHousingTypeName
* @Description 房屋类型名称
* @Parameter housing_type int
* @Return data string
 */
func GetHousingTypeName(housing_type int) (data string) {
	housingTypeMap := map[int]string{1: "普通公寓", 2: "公寓", 3: "一户建", 4: "别墅", 5: "民宿", 6: "简易旅馆", 0: "其他"}
	return housingTypeMap[housing_type]
}

/*
* @Title GetHuxingAlias
* @Description 户型别名
* @Parameter huxing string
* @Return data string
 */
func GetHuxingAlias(huxing string) (data string) {
	huxingMap := map[string]string{"1R": "1室", "1K": "1室", "1DK": "1室1厅", "1LDK": "1室2厅", "2K": "2室", "2DK": "2室1厅", "2LDK": "2室2厅", "3K": "3室", "3DK": "3室1厅", "3LDK": "3室2厅", "4K": "4室", "4DK": "4室1厅", "4LDK": "4室2厅", "5K以上": "5室"}
	return huxingMap[huxing]
}

// 公用-意见反馈
func (this *PublicModel) Public_Feedback(types int, contact, content string) (errMsg string) {
	sql := `INSERT INTO p_feedback(type, contact, content) VALUES(?,?,?)`
	_, err := db.Db.Exec(sql, types, contact, content)
	if err != nil {
		return "意见反馈失败"
	}
	return
}

type PublicContactReturn struct {
	CompanyName string `json:"company_name"`
	Adress      string `json:"adress"`
	UserName    string `json:"user_name"`
	Telephone   string `json:"telephone"`
	Fax         string `json:"fax"`
	Email       string `json:"email"`
}

// 公用-联系方式
func (this *PublicModel) Public_Contact(estateId int) (data *PublicContactReturn, errMsg string) {
	// 发布该房源的公司
	sql := `SELECT u.company_id
			FROM p_estate e
			LEFT JOIN p_user u ON u.id=e.user_id
			WHERE e.id=? AND e.is_del=0`
	row, err := db.Db.Query(sql, estateId)
	if err != nil {
		return data, "获取发布房源公司失败"
	}
	if len(row) == 0 {
		return data, "该房源不存在"
	}
	companyId, _ := strconv.Atoi(string(row[0]["company_id"]))

	// 发布该房源的公司主管
	sql = `SELECT id
		   FROM p_user
		   WHERE company_id=? AND user_type=1`
	row, err = db.Db.Query(sql, companyId)
	if err != nil {
		return data, "获取发布该房源的公司主管失败"
	}
	userId, _ := strconv.Atoi(string(row[0]["id"]))

	// 公司信息
	companyInfo, errMsg := userModel.GetCompanyInfo(companyId)
	if errMsg != "" {
		return data, errMsg
	}

	// 用户信息
	userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: userId})
	if errMsg != "" {
		return data, errMsg
	}

	return &PublicContactReturn{
		CompanyName: companyInfo.Name,
		Adress:      companyInfo.Adress,
		UserName:    userInfo.Name,
		Telephone:   userInfo.Telephone,
		Fax:         userInfo.Fax,
		Email:       userInfo.Email,
	}, ""
}
