package v1

import (
	"encoding/json"
	"estate/db"
	"estate/utils"
	"strconv"
	"time"
)

type BaseModel struct{}

type BaseDateListReturn struct {
	AddTime []string `json:"add_time"`
}

// 本部中介-日期列表
func (this *BaseModel) Base_DateList() (data *BaseDateListReturn, errMsg string) {
	data = new(BaseDateListReturn)

	// 日期列表
	sql := `SELECT deal_time
			FROM p_recommend
			WHERE is_pay=1
			ORDER BY deal_time DESC`
	rows, err := db.Db.Query(sql)
	if err != nil {
		return data, "获取日期列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	dealTimeMap := make(map[string]int)
	for _, value := range rows {
		t, _ := time.ParseInLocation("2006-01-02", string(value["deal_time"]), time.Local)
		dealTime := t.Format("2006-01")
		if _, ok := dealTimeMap[dealTime]; !ok {
			dealTimeMap[dealTime] = 1
			data.AddTime = append(data.AddTime, dealTime)
		}
	}
	return
}

type BaseSalesAchievementReturn struct {
	TotalPrice int                    `json:"total_price"`
	List       []SalesAchievementList `json:"list"`
	Pagenation Pagenation             `json:"pagenation"`
}

type SalesAchievementList struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Count      int    `json:"count"`
	TotalPrice int    `json:"total_price"`
}

// 本部中介-销售业绩
func (this *BaseModel) Base_SalesAchievement(addTime string, perPage, lastId int) (data *BaseSalesAchievementReturn, errMsg string) {
	data = new(BaseSalesAchievementReturn)
	data.Pagenation.LastId = -1

	// 时间
	addTime += "-01"
	t, _ := time.ParseInLocation("2006-01-02", addTime, time.Local)
	t = t.AddDate(0, 1, 0)
	endTime := t.Format("2006-01-02")

	// 分页
	var having string
	if perPage == 0 {
		perPage = 10
	}
	if lastId > 0 {
		// 获取当前分页的最后一条数据的时间
		sql := `SELECT SUM(e.price) total_price
				FROM p_recommend r
				LEFT JOIN p_estate e ON e.id=r.estate_id
				WHERE r.user_id=?`
		row, err := db.Db.Query(sql, lastId)
		if err != nil {
			return data, "获取销售业绩失败"
		}
		totalPrice := string(row[0]["total_price"])
		having = ` HAVING (SUM(e.price)<` + totalPrice + ` OR (SUM(e.price)=` + totalPrice + ` AND r.user_id>` + strconv.Itoa(lastId) + `))`
	}

	// 销售业绩列表
	sql := `SELECT r.user_id, u.name, COUNT(r.user_id) count, SUM(e.price) total_price
			FROM p_recommend r
			LEFT JOIN p_user u ON u.id=r.user_id
			LEFT JOIN p_estate e ON e.id=r.estate_id
			WHERE is_pay=1 AND deal_time>=? AND deal_time<=?
			GROUP BY r.user_id
			ORDER BY total_price DESC, user_id ` + having + ` LIMIT 0,?`
	rows, err := db.Db.Query(sql, addTime, endTime, perPage+1)
	if err != nil {
		return data, "获取销售业绩失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	var totalPrice float64
	for key, value := range rows {
		if key < perPage {
			totalPrice += utils.Str2float64(string(value["total_price"]))
			data.List = append(data.List, SalesAchievementList{
				Id:         utils.Str2int(string(value["user_id"])),
				Name:       string(value["name"]),
				Count:      utils.Str2int(string(value["count"])),
				TotalPrice: utils.Str2int(string(value["total_price"])),
			})
			lastId = utils.Str2int(string(value["user_id"]))
		} else {
			data.Pagenation.LastId = lastId
		}
	}
	data.TotalPrice = utils.Float64Toint(totalPrice)

	return
}

// 本部中介-中介费用统计列表
func (this *BaseModel) Base_SalesProfitList(addTime string, perPage, lastId int) (data *BaseSalesAchievementReturn, errMsg string) {
	data = new(BaseSalesAchievementReturn)
	data.Pagenation.LastId = -1

	// 时间
	addTime += "-01"
	t, _ := time.ParseInLocation("2006-01-02", addTime, time.Local)
	t = t.AddDate(0, 1, 0)
	endTime := t.Format("2006-01-02")

	// 分页
	var having string
	if perPage == 0 {
		perPage = 10
	}
	if lastId > 0 {
		// 获取当前分页的最后一条数据的时间
		sql := `SELECT SUM(e.agency_fee) total_price
				FROM p_recommend r
				LEFT JOIN p_estate e ON e.id=r.estate_id
				WHERE r.user_id=?`
		row, err := db.Db.Query(sql, lastId)
		if err != nil {
			return data, "获取中介费用失败"
		}
		totalPrice := string(row[0]["agency_fee"])
		having = ` HAVING (SUM(e.agency_fee)<` + totalPrice + ` OR (SUM(e.agency_fee)=` + totalPrice + ` AND r.user_id>` + strconv.Itoa(lastId) + `))`
	}

	// 中介费用列表
	sql := `SELECT r.user_id, u.name, COUNT(r.user_id) count, SUM(e.agency_fee) total_price
			FROM p_recommend r
			LEFT JOIN p_user u ON u.id=r.user_id
			LEFT JOIN p_estate e ON e.id=r.estate_id
			WHERE is_pay=1 AND deal_time>=? AND deal_time<=?
			GROUP BY r.user_id
			ORDER BY total_price DESC, user_id ` + having + ` LIMIT 0,?`
	rows, err := db.Db.Query(sql, addTime, endTime, perPage+1)
	if err != nil {
		return data, "获取中介费用失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	var totalPrice float64 // 总额
	for key, value := range rows {
		if key < perPage {
			totalPrice += utils.Str2float64(string(value["agency_fee"]))

			// 没有中介id的归为其他
			name := string(value["name"])
			if utils.Str2int(string(value["id"])) == 0 {
				name = "其他"
			}

			data.List = append(data.List, SalesAchievementList{
				Id:         utils.Str2int(string(value["user_id"])),
				Name:       name,
				Count:      utils.Str2int(string(value["count"])),
				TotalPrice: utils.Str2int(string(value["agency_fee"])),
			})
			lastId = utils.Str2int(string(value["user_id"]))
		} else {
			data.Pagenation.LastId = lastId
		}
	}
	data.TotalPrice = utils.Float64Toint(totalPrice)

	return
}

type BaseSalesProfitDetailReturn struct {
	List []SalesProfitDetailList `json:"list"`
}

type SalesProfitDetailList struct {
	EstateId   int    `json:"estate_id"`
	EstateCode string `json:"estate_code"`
	EstateName string `json:"estate_name"`
	Price      int    `json:"price"`
	DealTime   string `json:"deal_time"`
}

// 本部中介-中介费用统计详情
func (this *BaseModel) Base_SalesProfitDetail(id int) (data *BaseSalesProfitDetailReturn, errMsg string) {
	data = new(BaseSalesProfitDetailReturn)

	// 中介费用详情
	sql := `SELECT r.estate_id, r.deal_time, e.code, e.agency_fee, e.region_id, re.name regionName, re.type regionType, e.housing_type
			FROM p_recommend r
			LEFT JOIN p_estate e ON e.id=r.estate_id
			LEFT JOIN p_region re ON re.id=e.region_id
			WHERE r.user_id=?`
	rows, err := db.Db.Query(sql, id)
	if err != nil {
		return data, "获取中介费用详情失败"
	}
	for _, value := range rows {
		// 地区名称
		var regionName, areaName string
		regionType, _ := strconv.Atoi(string(value["regionType"]))
		if regionType == 1 { // 市
			regionName = string(value["regionName"])
		} else { // 区
			sql := `SELECT name FROM p_region WHERE id=?`
			row, err := db.Db.Query(sql, utils.Str2int(string(value["regionPId"])))
			if err != nil {
				return data, "获取地区失败"
			}
			regionName, areaName = string(row[0]["name"]), string(value["regionName"])
		}

		// 房屋类型
		housingType := map[int]string{1: "普通公寓", 2: "公寓", 3: "一户建", 4: "别墅", 5: "民宿", 6: "简易旅馆"}
		housingName := housingType[utils.Str2int(string(value["housing_type"]))]

		// 详情
		data.List = append(data.List, SalesProfitDetailList{
			EstateId:   utils.Str2int(string(value["estate_id"])),
			EstateCode: string(value["code"]),
			EstateName: regionName + areaName + housingName,
			Price:      utils.Str2int(string(value["agency_fee"])),
			DealTime:   string(value["deal_time"]),
		})
	}

	return
}

type BaseSalesProfitSettingDetailReturn struct {
	Price     int   `json:"price"`
	AgencyFee int   `json:"agency_fee"`
	Buyer     Buyer `json:"buyer"`
	Seller    Buyer `json:"Seller"`
}

type Buyer struct {
	ServiceFee int    `json:"service_fee"`
	FixedFee   int    `json:"fixed_fee"`
	ExciseFee  string `json:"excise_fee"`
}

// 本部中介-中介费用统计设置详情
func (this *BaseModel) Base_SalesProfitSettingDetail(estateId int) (data *BaseSalesProfitSettingDetailReturn, errMsg string) {
	data = new(BaseSalesProfitSettingDetailReturn)

	// 房源详情
	sql := `SELECT price, agency_fee
			FROM p_estate
			WHERE id=?`
	row, err := db.Db.Query(sql, estateId)
	if err != nil {
		return data, "获取房源详情失败"
	}
	if len(row) == 0 {
		return data, "该房源不存在"
	}
	data = &BaseSalesProfitSettingDetailReturn{
		Price:     utils.Str2int(string(row[0]["price"])),
		AgencyFee: utils.Str2int(string(row[0]["agency_fee"])),
	}

	// 中介费详情
	sql = `SELECT service_fee, fixed_fee, excise_fee, user_type
		   FROM base_agency_fee
		   WHERE estate_id=?`
	rows, err := db.Db.Query(sql, estateId)
	if err != nil {
		return data, "获取中介详情失败"
	}
	if len(rows) == 0 { // 还未设置中介费，取本部系统中介费
		// 本部基础信息
		baseInfo, errMsg := userModel.GetBaseInfo()
		if errMsg != "" {
			return data, errMsg
		}
		data.Buyer = Buyer{
			ServiceFee: baseInfo.ServiceFee,
			FixedFee:   baseInfo.FixedFee,
			ExciseFee:  baseInfo.ExciseFee,
		}
		return data, ""
	}
	for _, value := range rows {
		if utils.Str2int(string(value["user_type"])) == 0 { // 买家
			data.Buyer = Buyer{
				ServiceFee: utils.Str2int(string(value["service_fee"])),
				FixedFee:   utils.Str2int(string(value["fixed_fee"])),
				ExciseFee:  string(value["excise_fee"]),
			}
		} else { // 卖家
			data.Seller = Buyer{
				ServiceFee: utils.Str2int(string(value["service_fee"])),
				FixedFee:   utils.Str2int(string(value["fixed_fee"])),
				ExciseFee:  string(value["excise_fee"]),
			}
		}
	}

	return
}

// 本部中介-中介费用统计设置修改
func (this *BaseModel) Base_SalesProfitSettingModify(estateId int, agencyJson string) (errMsg string) {
	// 解析json
	agency := new(BaseSalesProfitSettingDetailReturn)
	if err := json.Unmarshal([]byte(agencyJson), agency); err != nil {
		return "解析json失败"
	}

	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 修改房源售价、中介费
	sql := `UPDATE p_estate
			SET price=?, agency_fee=?
			WHERE id=?`
	_, err := transaction.Exec(sql, estateId, agency.Price, agency.AgencyFee)
	if err != nil {
		transaction.Rollback()
		return "修改房源信息失败"
	}

	// 删除旧的中介信息
	sql = `DELETE FROM base_agency_fee WHERE estate_id=?`
	_, err = transaction.Exec(sql, estateId)
	if err != nil {
		transaction.Rollback()
		return "删除中介信息失败"
	}

	// 更新买家中介设置
	sql = `INSERT INTO base_agency_fee(estate_id, service_fee, fixed_fee, excise_fee, user_type) VALUES(?,?,?,?,?)`
	_, err = transaction.Exec(sql, estateId, agency.Buyer.ServiceFee, agency.Buyer.FixedFee, agency.Buyer.ExciseFee, 0)
	if err != nil {
		transaction.Rollback()
		return "更新买家中介设置失败"
	}

	// 更新卖家中介设置
	if agency.Seller.ServiceFee > 0 {
		sql = `INSERT INTO base_agency_fee(estate_id, service_fee, fixed_fee, excise_fee, user_type) VALUES(?,?,?,?,?)`
		_, err = transaction.Exec(sql, estateId, agency.Seller.ServiceFee, agency.Seller.FixedFee, agency.Seller.ExciseFee, 1)
		if err != nil {
			transaction.Rollback()
			return "更新卖家中介设置失败"
		}
	}

	// 提交事务
	if err = transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}

	return
}

type BaseWaitDistributionListReturn struct {
	List       []WaitDistributionList `json:"list"`
	Pagenation Pagenation             `json:"pagenation"`
}

type WaitDistributionList struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Sex     int    `json:"sex"`
	Wechat  string `json:"wechat"`
	Source  string `json:"source"`
	AddTime string `json:"add_time"`
}

// 本部中介-待分配客户列表
func (this *BaseModel) Base_WaitDistributionList(perPage, lastId int) (data *BaseWaitDistributionListReturn, errMsg string) {
	data = new(BaseWaitDistributionListReturn)
	data.Pagenation.LastId = -1

	// 分页
	var where string
	if perPage == 0 {
		perPage = 10
	}
	if lastId > 0 {
		// 获取当前分页的最后一条数据的时间
		sql := `SELECT add_time FROM p_recommend WHERE id=?`
		row, err := db.Db.Query(sql, lastId)
		if err != nil {
			return data, "获取推荐时间失败"
		}
		addTime := string(row[0]["add_time"])
		where += ` AND (r.add_time<"` + addTime + `" OR (r.add_time="` + addTime + `" AND r.id<` + strconv.Itoa(lastId) + `))`
	}

	// 待分配客户列表
	sql := `SELECT r.id, r.user_id, r.add_time, t.name, t.sex, t.wechat
			FROM p_recommend r
			LEFT JOIN p_tourists t ON t.id=r.tourists_id
			WHERE r.is_distribution=0 ` + where + ` ORDER BY r.add_time DESC, id DESC LIMIT 0,?`
	rows, err := db.Db.Query(sql, perPage+1)
	if err != nil {
		return data, "获取待分配客户失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for key, value := range rows {
		if key < perPage {
			var (
				company   string
				userId, _ = strconv.Atoi(string(value["user_id"]))
			)
			if userId > 0 {
				// 公司名称
				sql = `SELECT c.name
					   FROM p_company c
					   LEFT JOIN p_user u ON u.company_id=c.id
					   WHERE u.id=?`
				row, err := db.Db.Query(sql, userId)
				if err != nil {
					return data, "获取公司名称失败"
				}
				company = string(row[0]["name"])
			}

			data.List = append(data.List, WaitDistributionList{
				Id:      utils.Str2int(string(value["id"])),
				Name:    string(value["name"]),
				Sex:     utils.Str2int(string(value["sex"])),
				Wechat:  string(value["wechat"]),
				Source:  company,
				AddTime: string(value["add_time"]),
			})

			lastId = utils.Str2int(string(value["id"]))
		} else {
			data.Pagenation.LastId = lastId
		}
	}

	return
}

// 本部中介-待分配客户分配
func (this *BaseModel) Base_WaitDistributionDistribution(id, userId int) (errMsg string) {
	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 分配
	sql := `INSERT INTO base_distribution(recommend_id, user_id) VALUES(?,?)`
	_, err := transaction.Exec(sql, id, userId)
	if err != nil {
		transaction.Rollback()
		return "更新分配失败"
	}

	// 更新推荐状态
	sql = `UPDATE p_recommend
		   SET is_distribution=1
		   WHERE id=?`
	_, err = transaction.Exec(sql, id)
	if err != nil {
		transaction.Rollback()
		return "更新推荐状态失败"
	}

	// 提交事务
	if err = transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}

	return
}

// 本部中介-待分配客户删除
func (this *BaseModel) Base_WaitDistributionDel(id int) (errMsg string) {
	// 删除待分配客户
	sql := `DELETE r, t
			FROM p_recommend r
			LEFT JOIN p_tourists t ON t.id=r.tourists_id
			WHERE r.id=?`
	_, err := db.Db.Exec(sql, id)
	if err != nil {
		return "待分配客户删除失败"
	}

	return
}

type BaseJapanManageListReturn struct {
	List       []JapanManageList `json:"list"`
	Pagenation Pagenation        `json:"pagenation"`
}

type JapanManageList struct {
	Id          int    `json:"id"`
	CompanyName string `json:"company_name"`
	UserName    string `json:"user_name"`
	Telephone   string `json:"telephone"`
	ExpiryDate  string `json:"expiry_date"`
}

// 本部中介-日本中介管理列表
func (this *BaseModel) Base_JapanManageList(keyword string, status, perPage, lastId int) (data *BaseJapanManageListReturn, errMsg string) {
	data = new(BaseJapanManageListReturn)
	data.Pagenation.LastId = -1

	// 条件
	var where string

	// 关键字
	if keyword != "" {
		where = ` AND (c.name='` + keyword + `' OR c.adress='` + keyword + `')`
	}

	// 是否过期
	switch status {
	case 1: // 未过期
		where += ` AND expiry_date>=curdate() `
	case 2: // 已过期
		where += ` AND expiry_date<curdate() `
	}

	// 分页
	if perPage == 0 {
		perPage = 10
	}
	if lastId > 0 {
		// 获取当前分页的最后一条数据的时间
		sql := `SELECT add_time FROM p_company WHERE id=?`
		row, err := db.Db.Query(sql, lastId)
		if err != nil {
			return data, "获取会社新建时间失败"
		}
		addTime := string(row[0]["add_time"])
		where += ` AND (c.add_time<"` + addTime + `" OR (c.add_time="` + addTime + `" AND c.id<` + strconv.Itoa(lastId) + `))`
	}

	// 列表
	sql := `SELECT c.id, c.name companyName, c.expiry_date, u.name userName, u.telephone
			FROM p_company c
			LEFT JOIN p_user u ON u.company_id=c.id
			WHERE c.group_id=3 AND u.user_type=0 ` + where + ` ORDER BY c.add_time DESC, c.id DESC LIMIT 0,?`
	rows, err := db.Db.Query(sql, perPage+1)
	if err != nil {
		return data, "获取列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for key, value := range rows {
		if key < perPage {
			data.List = append(data.List, JapanManageList{
				Id:          utils.Str2int(string(value["id"])),
				CompanyName: string(value["companyName"]),
				UserName:    string(value["userName"]),
				Telephone:   string(value["telephone"]),
				ExpiryDate:  string(value["expiry_date"]),
			})
			lastId = utils.Str2int(string(value["id"]))
		} else {
			data.Pagenation.LastId = lastId
		}
	}

	return
}

type BaseJapanManageDetailReturn struct {
	Id          int    `json:"id"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	UserId      int    `json:"user_id"`
	UserName    string `json:"user_name"`
	Telephone   string `json:"telephone"`
	Fax         string `json:"fax"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ExpiryDate  string `json:"expiry_date"`
}

// 本部中介-日本中介管理详情
func (this *BaseModel) Base_JapanManageDetail(id int) (data *BaseJapanManageDetailReturn, errMsg string) {
	// 详情
	sql := `SELECT c.id, c.name companyName, c.adress, c.expiry_date, u.id userId, u.name userName, u.telephone, u.fax, u.email
			FROM p_company c
			LEFT JOIN p_user u ON u.company_id=c.id
			WHERE c.id=? AND u.user_type=0 AND c.group_id=3`
	row, err := db.Db.Query(sql, id)
	if err != nil {
		return nil, "获取详情失败"
	}
	if len(row) == 0 {
		return nil, ""
	}
	return &BaseJapanManageDetailReturn{
		Id:          utils.Str2int(string(row[0]["id"])),
		CompanyName: string(row[0]["companyName"]),
		Address:     string(row[0]["adress"]),
		UserId:      utils.Str2int(string(row[0]["user_id"])),
		UserName:    string(row[0]["userName"]),
		Telephone:   string(row[0]["telephone"]),
		Fax:         string(row[0]["fax"]),
		Email:       string(row[0]["email"]),
		ExpiryDate:  string(row[0]["expiry_date"]),
	}, ""
}

// 本部中介-日本中介管理添加/编辑
func (this *BaseModel) Base_JapanManageAdd(addParam *BaseJapanManageDetailReturn) (errMsg string) {
	// 判断除自己以外是否还存在该邮箱
	_, errMsg = publicModel.ExistEmail(addParam.Email, addParam.UserId)
	if errMsg != "" {
		return
	}

	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 根据id是否为0来判断是编辑还是新增
	if addParam.Id == 0 { // 新增
		// 更新公司
		sql := `INSERT INTO p_company(group_id, name, adress, expiry_date) VALUES(?,?,?,?)`
		row, err := transaction.Exec(sql, 3, addParam.CompanyName, addParam.Address, addParam.ExpiryDate)
		if err != nil {
			transaction.Rollback()
			return "更新公司失败"
		}
		lastId, _ := row.LastInsertId()

		// 更新公司主管
		sql = `INSERT INTO p_user(company_id, user_type, email, password, name, telephone, fax) VALUES(?,?,?,?,?,?,?)`
		_, err = transaction.Exec(sql, int(lastId), 1, addParam.Email, string(utils.HashPassword(addParam.Password)), addParam.UserName, addParam.Telephone, addParam.Fax)
		if err != nil {
			transaction.Rollback()
			return "更新公司主管失败"
		}
	} else { // 编辑
		// 更新公司
		sql := `UPDATE p_company
				SET name=?, adress=?, expiry_date=?
				WHERE id=?`
		_, err := transaction.Exec(sql, addParam.CompanyName, addParam.Address, addParam.ExpiryDate, addParam.Id)
		if err != nil {
			transaction.Rollback()
			return "更新公司失败."
		}

		// 更新公司主管
		var passwordSql string
		if addParam.Password != "" {
			passwordSql = ` ,password="` + string(utils.HashPassword(addParam.Password)) + `"`
		}
		sql = `UPDATE p_user
			   SET email=?, name=?, telephone=?, fax=? ` + passwordSql +
			`WHERE company_id=? AND user_type=1`
		_, err = transaction.Exec(sql, addParam.Email, addParam.UserName, addParam.Telephone, addParam.Fax, addParam.Id)
		if err != nil {
			transaction.Rollback()
			return "更新公司主管失败."
		}
	}

	// 提交事务
	if err := transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}

	return
}

// 本部中介-日本中介管理删除
func (this *BaseModel) Base_JapanManageDel(id int) (errMsg string) {
	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 房源软删除
	sql := `UPDATE p_estate e
			LEFT JOIN p_user u ON u.id=e.user_id
			SET e.is_del=1
			WHERE u.company_id=?`
	_, err := transaction.Exec(sql, id)
	if err != nil {
		transaction.Rollback()
		return "房源删除失败"
	}

	// 公司及员工硬删除
	sql = `DELETE c, u
		   FROM p_company c
		   LEFT JOIN p_user u ON u.company_id=c.id
		   WHERE c.id=?`
	_, err = transaction.Exec(sql, id)
	if err != nil {
		transaction.Rollback()
		return "公司及员工删除失败"
	}

	// 提交事务
	if err = transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}

	return
}

type BaseChinaManageRegionListReturn struct {
	List []ChinaManageRegionList `json:"list"`
}

type ChinaManageRegionList struct {
	RegionId   int    `json:"region_id"`
	RegionName string `json:"region_name"`
	Mark       int    `json:"mark"`
}

// 本部中介-中国中介管理地区列表
func (this *BaseModel) Base_ChinaManageRegionList() (data *BaseChinaManageRegionListReturn, errMsg string) {
	data = new(BaseChinaManageRegionListReturn)

	// 地区公司
	sql := `SELECT id, region_id 
			FROM p_company
			WHERE group_id=2`
	rows, err := db.Db.Query(sql)
	if err != nil {
		return data, "获取地区公司失败"
	}
	regionCompanyMap := make(map[int]int)
	if len(rows) > 0 {
		for _, value := range rows {
			regionCompanyMap[utils.Str2int(string(value["region_id"]))] = utils.Str2int(string(value["id"]))
		}
	}

	// 地区列表
	sql = `SELECT id, name
		   FROM p_region
		   WHERE group_id=2`
	rows, err = db.Db.Query(sql)
	if err != nil {
		return data, "获取地区列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for _, value := range rows {
		var (
			mark        int
			regionId, _ = strconv.Atoi(string(value["id"]))
		)
		if _, ok := regionCompanyMap[regionId]; ok {
			mark = 1
		}
		data.List = append(data.List, ChinaManageRegionList{
			RegionId:   regionId,
			RegionName: string(value["name"]),
			Mark:       mark,
		})
	}

	return
}

type BaseChinaManageListReturn struct {
	List       []ChinaManageList `json:"list"`
	Pagenation Pagenation        `json:"pagenation"`
}

type ChinaManageList struct {
	Id          int    `json:"id"`
	CompanyName string `json:"company_name"`
	UserName    string `json:"user_name"`
	Telephone   string `json:"telephone"`
	AddTime     string `json:"add_time"`
}

// 本部中介-中国中介管理列表
func (this *BaseModel) Base_ChinaManageList(keyword string, regionId, perPage, lastId int) (data *BaseChinaManageListReturn, errMsg string) {
	data = new(BaseChinaManageListReturn)
	data.Pagenation.LastId = -1

	// 条件
	var where string

	// 关键字
	if keyword != "" {
		where = ` AND (c.name='` + keyword + `' OR c.adress='` + keyword + `')`
	}

	// 地区id
	if regionId > 0 {
		where += ` AND c.region_id=` + strconv.Itoa(regionId)
	}

	// 分页
	if perPage == 0 {
		perPage = 10
	}
	if lastId > 0 {
		// 获取当前分页的最后一条数据的时间
		sql := `SELECT add_time FROM p_company WHERE id=?`
		row, err := db.Db.Query(sql, lastId)
		if err != nil {
			return data, "获取公司新建时间失败"
		}
		addTime := string(row[0]["add_time"])
		where += ` AND (c.add_time<"` + addTime + `" OR (c.add_time="` + addTime + `" AND c.id<` + strconv.Itoa(lastId) + `))`
	}

	// 列表
	sql := `SELECT c.id, c.name companyName, c.add_time, u.name userName, u.telephone
			FROM p_company c
			LEFT JOIN p_user u ON u.company_id=c.id
			WHERE c.group_id=2 AND u.user_type=0 ` + where + ` ORDER BY c.add_time DESC, c.id DESC LIMIT 0,?`
	rows, err := db.Db.Query(sql, perPage+1)
	if err != nil {
		return data, "获取列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for key, value := range rows {
		if key < perPage {
			data.List = append(data.List, ChinaManageList{
				Id:          utils.Str2int(string(value["id"])),
				CompanyName: string(value["companyName"]),
				UserName:    string(value["userName"]),
				Telephone:   string(value["telephone"]),
				AddTime:     string(value["add_time"]),
			})
			lastId = utils.Str2int(string(value["id"]))
		} else {
			data.Pagenation.LastId = lastId
		}
	}

	return
}

type BaseChinaManageDetailReturn struct {
	Id          int    `json:"id"`
	RegionId    int    `json:"region_id"`
	RegionName  string `json:"region_name"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	UserId      int    `json:"user_id"`
	UserName    string `json:"user_name"`
	Telephone   string `json:"telephone"`
	Fax         string `json:"fax"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// 本部中介-中国中介管理详情
func (this *BaseModel) Base_ChinaManageDetail(id int) (data *BaseChinaManageDetailReturn, errMsg string) {
	// 详情
	sql := `SELECT c.id, c.region_id, r.name regionName, c.name companyName, c.adress, u.id userId, u.name userName, u.telephone, u.fax, u.email
			FROM p_company c
			LEFT JOIN p_region r ON r.id=c.region_id
			LEFT JOIN p_user u ON u.company_id=c.id
			WHERE c.id=? AND u.user_type=0 AND c.group_id=2`
	row, err := db.Db.Query(sql, id)
	if err != nil {
		return nil, "获取详情失败"
	}
	if len(row) == 0 {
		return nil, ""
	}
	return &BaseChinaManageDetailReturn{
		Id:          utils.Str2int(string(row[0]["id"])),
		RegionId:    utils.Str2int(string(row[0]["region_id"])),
		RegionName:  string(row[0]["regionName"]),
		CompanyName: string(row[0]["companyName"]),
		Address:     string(row[0]["adress"]),
		UserId:      utils.Str2int(string(row[0]["user_id"])),
		UserName:    string(row[0]["userName"]),
		Telephone:   string(row[0]["telephone"]),
		Fax:         string(row[0]["fax"]),
		Email:       string(row[0]["email"]),
	}, ""
}

// 本部中介-中国中介管理添加/编辑
func (this *BaseModel) Base_ChinaManageAdd(addParam *BaseChinaManageDetailReturn) (errMsg string) {
	// 判断除自己以外是否还存在该邮箱
	_, errMsg = publicModel.ExistEmail(addParam.Email, addParam.UserId)
	if errMsg != "" {
		return
	}

	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 根据id是否为0来判断是编辑还是新增
	if addParam.Id == 0 { // 新增
		// 更新公司
		sql := `INSERT INTO p_company(group_id, region_id, name, adress ) VALUES(?,?,?,?)`
		row, err := transaction.Exec(sql, 2, addParam.RegionId, addParam.CompanyName, addParam.Address)
		if err != nil {
			transaction.Rollback()
			return "更新公司失败"
		}
		lastId, _ := row.LastInsertId()

		// 更新公司主管
		sql = `INSERT INTO p_user(company_id, user_type, email, password, name, telephone, fax) VALUES(?,?,?,?,?,?,?)`
		_, err = transaction.Exec(sql, int(lastId), 1, addParam.Email, string(utils.HashPassword(addParam.Password)), addParam.UserName, addParam.Telephone, addParam.Fax)
		if err != nil {
			transaction.Rollback()
			return "更新公司主管失败"
		}
	} else { // 编辑
		// 更新公司
		sql := `UPDATE p_company
				SET region_id=?, name=?, adress=?
				WHERE id=?`
		_, err := transaction.Exec(sql, addParam.RegionId, addParam.CompanyName, addParam.Address, addParam.Id)
		if err != nil {
			transaction.Rollback()
			return "更新公司失败."
		}

		// 更新公司主管
		var passwordSql string
		if addParam.Password != "" {
			passwordSql = ` ,password="` + string(utils.HashPassword(addParam.Password)) + `"`
		}
		sql = `UPDATE p_user
			   SET email=?, name=?, telephone=?, fax=? ` + passwordSql +
			`WHERE company_id=? AND user_type=1`
		_, err = transaction.Exec(sql, addParam.Email, addParam.UserName, addParam.Telephone, addParam.Fax, addParam.Id)
		if err != nil {
			transaction.Rollback()
			return "更新公司主管失败."
		}
	}

	// 提交事务
	if err := transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}

	return
}

// 本部中介-中国中介管理删除
func (this *BaseModel) Base_ChinaManageDel(id int) (errMsg string) {
	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 推荐关系删除
	sql := `UPDATE p_recommend r
			LEFT JOIN p_user u ON u.id=r.user_id
			SET r.user_id=0
			WHERE u.company_id=?`
	_, err := transaction.Exec(sql, id)
	if err != nil {
		transaction.Rollback()
		return "推荐关系删除失败"
	}

	// 公司及员工硬删除
	sql = `DELETE c, u
		   FROM p_company c
		   LEFT JOIN p_user u ON u.company_id=c.id
		   WHERE c.id=?`
	_, err = transaction.Exec(sql, id)
	if err != nil {
		transaction.Rollback()
		return "公司及员工删除失败"
	}

	// 提交事务
	if err = transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}

	return
}

type BaseCustomerManageSourceListReturn struct {
	List []CustomerManageSourceList `json:"list"`
}

type CustomerManageSourceList struct {
	CompanyId   int    `json:"company_id"`
	CompanyName string `json:"company_name"`
}

// 本部中介-客户管理来源列表
func (this *BaseModel) Base_CustomerManageSourceList() (data *BaseCustomerManageSourceListReturn, errMsg string) {
	data = new(BaseCustomerManageSourceListReturn)

	// 列表
	sql := `SELECT u.company_id, c.name
			FROM p_recommend r
			LEFT JOIN p_user u ON u.id=r.user_id
			LEFT JOIN p_company c ON c.id=u.company_id
			WHERE r.is_distribution=1
			GROUP BY u.company_id`
	rows, err := db.Db.Query(sql)
	if err != nil {
		return data, "获取列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for _, value := range rows {
		var (
			companyId, _ = strconv.Atoi(string(value["company_id"]))
			companyName  = string(value["name"])
		)
		if companyId == 0 {
			companyId, companyName = -1, "游客咨询"
		}
		data.List = append(data.List, CustomerManageSourceList{
			CompanyId:   companyId,
			CompanyName: companyName,
		})
	}
	return
}

type BaseCustomerManageListParamater struct {
	Keyword   string
	UserId    int
	CompanyId int
	Status    int
	PerPage   int
	LastId    int
}

type BaseCustomerManageListReturn struct {
	List       []CustomerManageList `json:"list"`
	Pagenation Pagenation           `json:"pagenation"`
}

type CustomerManageList struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Sex       int    `json:"sex"`
	Wechat    string `json:"wechat"`
	AddTime   string `json:"add_time"`
	IsButt    int    `json:"is_butt"`
	IsToJapan int    `json:"is_to_japan"`
	IsAgree   int    `json:"is_agree"`
	IsPay     int    `json:"is_pay"`
	IsLoan    int    `json:"is_loan"`
}

// 本部中介-客户管理列表
func (this *BaseModel) Base_CustomerManageList(cusParam *BaseCustomerManageListParamater) (data *BaseCustomerManageListReturn, errMsg string) {
	data = new(BaseCustomerManageListReturn)
	data.Pagenation.LastId = -1

	// 条件
	var where string

	// 关键字
	if cusParam.Keyword != "" {
		where = ` AND (t.name="` + cusParam.Keyword + `" OR t.wechat="` + cusParam.Keyword + `" OR e.code="` + cusParam.Keyword + `")`
	}

	// 销售id
	if cusParam.UserId > 0 {
		where += ` AND d.user_id=` + strconv.Itoa(cusParam.UserId)
	}

	// 公司id
	if cusParam.CompanyId == -1 { // 游客咨询
		where += ` AND r.user_id=0`
	} else if cusParam.CompanyId > 0 { // 公司推荐
		where += ` AND u.company_id=` + strconv.Itoa(cusParam.CompanyId)
	}

	// 状态
	switch cusParam.Status {
	case 1:
		where += ` AND r.is_butt=0`
	case 2:
		where += ` AND r.is_butt=1`
	case 3:
		where += ` AND r.is_to_japan=0`
	case 4:
		where += ` AND r.is_to_japan=1`
	case 5:
		where += ` AND r.is_agree=0`
	case 6:
		where += ` AND r.is_agree=1`
	case 7:
		where += ` AND r.is_pay=0`
	case 8:
		where += ` AND r.is_pay=1`
	case 9:
		where += ` AND r.is_loan=0`
	case 10:
		where += ` AND r.is_loan=1`
	default:
		return data, "不存在该状态"
	}

	// 分页
	if cusParam.PerPage == 0 {
		cusParam.PerPage = 10
	}
	if cusParam.LastId > 0 {
		// 获取当前分页的最后一条数据的时间
		sql := `SELECT add_time FROM p_recommend WHERE id=?`
		row, err := db.Db.Query(sql, cusParam.LastId)
		if err != nil {
			return data, "获取房源推荐时间失败"
		}
		addTime := string(row[0]["add_time"])
		where += ` AND (r.add_time<"` + addTime + `" OR (e.add_time="` + addTime + `" AND r.id<` + strconv.Itoa(cusParam.LastId) + `))`
	}

	// 列表
	sql := `SELECT r.id, r.add_time, r.is_butt, r.is_to_japan, r.is_agree, r.is_pay, r.is_loan, t.name, t.sex, t.wechat
			FROM p_recommend r
			LEFT JOIN p_tourists t ON t.id=r.tourists_id
			LEFT JOIN base_distribution d ON d.recommend_id=r.id
			LEFT JOIN p_user u ON u.id=r.user_id
			LEFT JOIN p_estate e ON e.id=r.estate_id
			WHERE r.is_distribution=1 ` + where + ` ORDER BY r.add_time DESC, id DESC LIMIT 0,?`
	rows, err := db.Db.Query(sql, cusParam.PerPage+1)
	if err != nil {
		return data, "获取客户列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for key, value := range rows {
		if key < cusParam.PerPage {
			data.List = append(data.List, CustomerManageList{
				Id:        utils.Str2int(string(value["id"])),
				Name:      string(value["name"]),
				Sex:       utils.Str2int(string(value["sex"])),
				Wechat:    string(value["wechat"]),
				AddTime:   string(value["add_time"]),
				IsButt:    utils.Str2int(string(value["is_butt"])),
				IsToJapan: utils.Str2int(string(value["is_to_japan"])),
				IsAgree:   utils.Str2int(string(value["is_agree"])),
				IsPay:     utils.Str2int(string(value["is_pay"])),
				IsLoan:    utils.Str2int(string(value["is_loan"])),
			})

			cusParam.LastId = utils.Str2int(string(value["id"]))
		} else {
			data.Pagenation.LastId = cusParam.LastId
		}
	}

	return
}

type BaseCustomerManageDetailReturn struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Sex        int    `json:"sex"`
	Wechat     string `json:"wechat"`
	IsButt     int    `json:"is_butt"`
	IsToJapan  int    `json:"is_to_japan"`
	IsAgree    int    `json:"is_agree"`
	IsPay      int    `json:"is_pay"`
	IsLoan     int    `json:"is_loan"`
	EstateCode string `json:"estate_code"`
	Price      string `json:"price"`
}

// 本部中介-客户管理详情
func (this *BaseModel) Base_CustomerManageDetail(id int) (data *BaseCustomerManageDetailReturn, errMsg string) {
	// 详情
	sql := `SELECT r.id, r.is_butt, r.is_to_japan, r.is_agree, r.is_pay, r.is_loan, t.name, t.sex, t.wechat, e.code, e.price
			FROM p_recommend r
			LEFT JOIN p_tourists t ON t.id=r.tourists_id
			LEFT JOIN p_estate e ON e.id=r.estate_id
			WHERE r.id=?`
	row, err := db.Db.Query(sql, id)
	if err != nil {
		return data, "获取客户详情失败"
	}
	return &BaseCustomerManageDetailReturn{
		Id:         utils.Str2int(string(row[0]["id"])),
		Name:       string(row[0]["name"]),
		Sex:        utils.Str2int(string(row[0]["sex"])),
		Wechat:     string(row[0]["wechat"]),
		IsButt:     utils.Str2int(string(row[0]["is_butt"])),
		IsToJapan:  utils.Str2int(string(row[0]["is_to_japan"])),
		IsAgree:    utils.Str2int(string(row[0]["is_agree"])),
		IsPay:      utils.Str2int(string(row[0]["is_pay"])),
		IsLoan:     utils.Str2int(string(row[0]["is_loan"])),
		EstateCode: string(row[0]["code"]),
		Price:      string(row[0]["price"]),
	}, ""
}

// 本部中介-客户管理编辑
func (this *BaseModel) Base_CustomerManageEdit(cusParam *BaseCustomerManageDetailReturn) (errMsg string) {
	// 更新客户信息
	sql := `UPDATE p_recommend r
			LEFT JOIN p_tourists t ON t.id=r.tourists_id
			LEFT JOIN p_estate e ON e.id=r.estate_id
			SET r.is_butt=?, r.is_to_japan=?, r.is_agree=?, r.is_pay=?, r.is_loan=?, t.name=?, t.sex=?, t.wechat=?, e.code=?, e.price=?
			WHERE r.id=?`
	_, err := db.Db.Exec(sql, cusParam.IsButt, cusParam.IsToJapan, cusParam.IsAgree, cusParam.IsPay, cusParam.IsLoan, cusParam.Name, cusParam.Sex, cusParam.Wechat, cusParam.EstateCode, cusParam.Price)
	if err != nil {
		return "更新客户信息失败"
	}
	return
}

// 本部中介-客户管理删除
func (this *BaseModel) Base_CustomerManageDel(id int) (errMsg string) {
	// 获取推荐状态
	sql := `SELECT is_pay
			FROM p_recommend
			WHERE id=?`
	row, err := db.Db.Query(sql, id)
	if err != nil {
		return "获取推荐状态失败"
	}
	if len(row) == 0 {
		return "该客户记录不存在"
	}
	isPay, _ := strconv.Atoi(string(row[0]["is_pay"]))
	if isPay == 1 {
		return "该客户已经付款，不允许删除"
	}

	// 删除
	sql = `DELETE d, r
		   FROM p_recommend r
		   LEFT JOIN base_distribution d ON d.recommend_id=r.id
		   WHERE r.id=?`
	_, err = db.Db.Exec(sql, id)
	if err != nil {
		return "删除客户失败"
	}

	return
}

type BaseProtectionPeriodShowReturn struct {
	ProtectionPeriod int `json:"protection_period"`
}

// 本部中介-保护期显示
func (this *BaseModel) Base_ProtectionPeriodShow() (data *BaseProtectionPeriodShowReturn, errMsg string) {
	// 本部基础信息
	baseInfo, errMsg := userModel.GetBaseInfo()
	if errMsg != "" {
		return data, errMsg
	}
	return &BaseProtectionPeriodShowReturn{
		ProtectionPeriod: baseInfo.ProtectionPeriod,
	}, ""
}

// 本部中介-保护期设置
func (this *BaseModel) Base_ProtectionPeriodSet(protectionPeriod int) (errMsg string) {
	// 更新保护期
	sql := `UPDATE base_info
			SET protection_period=?
			WHERE id=1`
	_, err := db.Db.Exec(sql, protectionPeriod)
	if err != nil {
		return "更新保护期失败"
	}
	return
}

type BaseAgencyFeeShowReturn struct {
	ServiceFee int    `json:"service_fee"`
	FixedFee   int    `json:"fixed_fee"`
	ExciseFee  string `json:"excise_fee"`
}

// 本部中介-中介费显示
func (this *BaseModel) Base_AgencyFeeShow() (data *BaseAgencyFeeShowReturn, errMsg string) {
	// 本部基础信息
	baseInfo, errMsg := userModel.GetBaseInfo()
	if errMsg != "" {
		return data, errMsg
	}
	return &BaseAgencyFeeShowReturn{
		ServiceFee: baseInfo.ServiceFee,
		FixedFee:   baseInfo.FixedFee,
		ExciseFee:  baseInfo.ExciseFee,
	}, ""
}

// 本部中介-中介费设置
func (this *BaseModel) Base_AgencyFeeSet(serviceFee, fixedFee int, exciseFee string) (errMsg string) {
	// 更新中介费
	sql := `UPDATE base_info
			SET service_fee=?, fixed_fee=?, excise_fee=?
			WHERE id=1`
	_, err := db.Db.Exec(sql, serviceFee, fixedFee, exciseFee)
	if err != nil {
		return "更新中介费失败"
	}
	return
}

type BaseNotifySetReturn struct {
	IsNotified int `json:"is_notified"`
}

// 本部中介-通知设置
func (this *BaseModel) Base_NotifySet() (data *BaseNotifySetReturn, errMsg string) {
	// 本部基础信息
	baseInfo, errMsg := userModel.GetBaseInfo()
	if errMsg != "" {
		return data, errMsg
	}

	// 根据当前通知状态来获取需要更新的状态
	var isNotified int
	if baseInfo.IsNotified == 0 {
		isNotified = 1
	}

	// 更新通知状态
	sql := `UPDATE base_info
			SET is_notified=?
			WHERE id=1`
	_, err := db.Db.Exec(sql, isNotified)
	if err != nil {
		return nil, "更新中介费失败"
	}

	return &BaseNotifySetReturn{
		IsNotified: isNotified,
	}, ""
}
