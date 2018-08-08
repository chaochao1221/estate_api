package v1

import (
	"estate/db"
	"estate/utils"
	"fmt"
	"strconv"
)

type ChinaModel struct{}

// 公用-推荐
func (this *ChinaModel) China_Recommend(estateId, sex, userId, groupId int, name, wechat string) (errMsg string) {
	// 判断该房源是否存在
	sql := `SELECT id
			FROM p_estate
			WHERE id=? AND is_del=0 AND status=1`
	rows, err := db.Db.Query(sql, estateId)
	if err != nil {
		return "获取房源失败"
	}
	if len(rows) == 0 {
		return "The estate does not exist"
	}

	// 用户信息
	userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: userId})
	if err != nil {
		return errMsg
	}

	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 更新游客
	sql = `INSERT INTO p_tourists(name, wechat, sex) VALUES(?,?,?)`
	row, err := transaction.Exec(sql, name, wechat, sex)
	if err != nil {
		transaction.Rollback()
		return "更新游客失败"
	}
	lastId, err := row.LastInsertId()

	// 更新推荐关系
	sql = `INSERT INTO p_recommend(estate_id, user_id, tourists_id) VALUES(?,?,?)`
	_, err = transaction.Exec(sql, estateId, userId, int(lastId))
	if err != nil {
		transaction.Rollback()
		return "更新推荐关系失败"
	}

	// 更新公司推荐数量
	sql = `UPDATE p_company
		   SET recommend_number=recommend_number+1
		   WHERE id=?`
	_, err = transaction.Exec(sql, userInfo.CompanyId)
	if err != nil {
		transaction.Rollback()
		return "更新公司推荐数量失败"
	}

	// 提交事务
	if err = transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}

	return
}

type ChinaCustomerProgressListParameter struct {
	Keyword     string
	IsButt      string
	IsToJapan   string
	IsAgree     string
	IsPay       string
	IsLoan      string
	UserId      int
	PerPage     int
	LastId      int
	LoginUserId int
	UserType    int
}

type ChinaCustomerProgressListReturn struct {
	List       []CustomerProgressList `json:"list"`
	Pagenation Pagenation             `json:"pagenation"`
}

type CustomerProgressList struct {
	Id        int    `json:"id"`
	EstateId  int    `json:"estate_id"`
	Name      string `json:"name"`
	Wechat    string `json:"wechat"`
	Sex       int    `json:"sex"`
	AddTime   string `json:"add_time"`
	IsButt    int    `json:"is_butt"`
	IsToJapan int    `json:"is_to_japan"`
	IsAgree   int    `json:"is_agree"`
	IsPay     int    `json:"is_pay"`
	IsLoan    int    `json:"is_loan"`
}

// 公用-客户进展列表
func (this *ChinaModel) China_CustomerProgressList(cusParam *ChinaCustomerProgressListParameter) (data *ChinaCustomerProgressListReturn, errMsg string) {
	data = new(ChinaCustomerProgressListReturn)
	data.Pagenation.LastId = -1

	// 条件
	var where string

	// 搜索关键字
	if cusParam.Keyword != "" {
		where += ` AND (t.name="` + cusParam.Keyword + `" OR t.wechat="` + cusParam.Keyword + `" OR e.code="` + cusParam.Keyword + `")`
	}

	// 主管可以看指定销售推荐的客户，销售只能看自己推荐的客户
	if cusParam.UserType == 0 { // 销售
		where += ` AND r.user_id=` + strconv.Itoa(cusParam.LoginUserId)
	} else { // 主管
		if cusParam.UserId > 0 {
			where += ` AND r.user_id=` + strconv.Itoa(cusParam.UserId)
		}
	}

	// 筛选状态
	if cusParam.IsButt != "" {
		where += ` AND r.is_butt=` + cusParam.IsButt
	}
	if cusParam.IsToJapan != "" {
		where += ` AND r.is_to_japan=` + cusParam.IsToJapan
	}
	if cusParam.IsAgree != "" {
		where += ` AND r.is_agree=` + cusParam.IsAgree
	}
	if cusParam.IsPay != "" {
		where += ` AND r.is_pay=` + cusParam.IsPay
	}
	if cusParam.IsLoan != "" {
		where += ` AND r.is_loan=` + cusParam.IsLoan
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
			return data, "获取房源发布时间失败"
		}
		addTime := string(row[0]["add_time"])
		where += ` AND (r.add_time<"` + addTime + `" OR (r.add_time="` + addTime + `" AND r.id<` + strconv.Itoa(cusParam.LastId) + `))`
	}

	// 获取推荐客户列表
	sql := `SELECT r.id, r.estate_id, r.add_time, r.is_butt, r.is_to_japan, r.is_agree, r.is_pay, r.is_loan, t.name, t.wechat, t.sex
			FROM p_recommend r
			LEFT JOIN p_estate e ON e.id=r.estate_id
			LEFT JOIN p_tourists t ON t.id=r.tourists_id
			WHERE r.user_id>0 ` + where + ` ORDER BY r.add_time DESC, id DESC LIMIT 0,?`
	rows, err := db.Db.Query(sql, cusParam.PerPage+1)
	if err != nil {
		fmt.Println(err, sql)
		return nil, "获取客户进展列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for key, value := range rows {
		if key < cusParam.PerPage {
			data.List = append(data.List, CustomerProgressList{
				Id:        utils.Str2int(string(value["id"])),
				EstateId:  utils.Str2int(string(value["estate_id"])),
				Name:      string(value["name"]),
				Wechat:    string(value["wechat"]),
				Sex:       utils.Str2int(string(value["sex"])),
				AddTime:   DateFormat(string(value["add_time"]), "2006/01/02"),
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

// 中国中介-客户进展删除
func (this *ChinaModel) China_CustomerProgressDel(id int) (errMsg string) {
	sql := `UPDATE p_recommend SET user_id=0 WHERE id=?`
	_, err := db.Db.Exec(sql, id)
	if err != nil {
		return "客户进展删除失败"
	}
	return
}
