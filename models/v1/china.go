package v1

import (
	"estate/db"
	"estate/utils"
	"strconv"
)

type ChinaModel struct{}

// 公用-推荐
func (this *ChinaModel) China_Recommend(estateId, sex, userId, groupId int, name, wechat string) (errMsg string) {
	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 更新游客
	sql := `INSERT INTO p_tourists(name, wechat, sex) VALUES(?,?,?)`
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

	// 提交事务
	if err = transaction.Commit(); err != nil {
		transaction.Rollback()
		return "提交事务失败"
	}

	return
}

type ChinaCustomerProgressListReturn struct {
	List       []CustomerProgressList `json:"list"`
	Pagenation Pagenation             `json:"pagenation"`
}

type CustomerProgressList struct {
	Id       int    `json:"id"`
	EstateId int    `json:"estate_id"`
	Name     string `json:"name"`
	Wechat   string `json:"wechat"`
	Sex      int    `json:"sex"`
	Status   int    `json:"status"`
	AddTime  string `json:"add_time"`
}

// 公用-客户进展列表
func (this *ChinaModel) China_CustomerProgressList(keyword string, status, userId, perPage, lastId, uId, userType int) (data *ChinaCustomerProgressListReturn, errMsg string) {
	data = new(ChinaCustomerProgressListReturn)
	data.Pagenation.LastId = -1

	// 主管可以看指定销售推荐的客户，销售只能看自己推荐的客户
	// 筛选状态（1:未对接 2:对接中 3:未赴日 4:已赴日 5:未成约 6:已成约 7:未付款 8:已付款 9:无贷款 10:已贷款）
	// 搜索关键字
	var where string
	if userType == 0 { // 销售
		where = ` AND r.user_id=` + strconv.Itoa(uId)
	} else { // 主管
		if userId > 0 {
			where = ` AND r.user_id=` + strconv.Itoa(userId)
		}
	}
	if status > 0 {
		where += ` AND r.status=` + strconv.Itoa(status)
	}
	if keyword != "" {
		where += ` AND (t.name="` + keyword + `" OR t.wechat="` + keyword + `" OR e.code="` + keyword + `")`
	}

	// 分页
	if perPage == 0 {
		perPage = 10
	}
	if lastId > 0 {
		// 获取当前分页的最后一条数据的时间
		sql := `SELECT add_time FROM p_recommend WHERE id=?`
		row, err := db.Db.Query(sql, lastId)
		if err != nil {
			return data, "获取房源发布时间失败"
		}
		addTime := string(row[0]["add_time"])
		where += ` AND (r.add_time<"` + addTime + `" OR (e.add_time="` + addTime + `" AND r.id<` + strconv.Itoa(lastId) + `))`
	}

	// 获取推荐客户列表
	sql := `SELECT r.id, r.estate_id, r.status, r.add_time, t.name, t.wechat, t.sex
			FROM p_recommend r
			LEFT JOIN p_estate e ON e.id=r.estate_id
			LEFT JOIN p_tourists t ON t.id=r.tourists_id
			WHERE r.user_id>0 ` + where + ` ORDER BY r.add_time DESC, id DESC LIMIT 0,?`
	rows, err := db.Db.Query(sql, perPage+1)
	if err != nil {
		return nil, "获取客户进展列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for key, value := range rows {
		if key < perPage {
			data.List = append(data.List, CustomerProgressList{
				Id:       utils.Str2int(string(value["id"])),
				EstateId: utils.Str2int(string(value["estate_id"])),
				Name:     string(value["name"]),
				Wechat:   string(value["wechat"]),
				Sex:      utils.Str2int(string(value["sex"])),
				Status:   utils.Str2int(string(value["status"])),
				AddTime:  string(value["add_time"]),
			})
			lastId = utils.Str2int(string(value["id"]))
		} else {
			data.Pagenation.LastId = lastId
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
