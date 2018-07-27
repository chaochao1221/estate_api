package v1

import (
	"estate/db"
	"estate/utils"
)

type PublicModel struct{}

var userModel = new(UserModel)

type PublicCompanyDetailReturn struct {
	GroupId          int    `json:"group_id"`
	UserId           int    `json:"user_id"`
	UserType         int    `json:"user_type"`
	CompanyName      string `json:"company_name"`
	CRecommendNumber int    `json:"c_recommend_number"`
	CDealNumber      int    `json:"c_deal_number"`
	EReleaseNumber   int    `json:"e_release_number"`
	EDealNumber      int    `json:"e_deal_number"`
	ExpiryDate       string `json:"expiry_date"`
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
		GroupId:          companyInfo.GroupId,
		UserId:           userId,
		UserType:         userInfo.UserType,
		CompanyName:      companyInfo.Name,
		CRecommendNumber: companyInfo.CRecommendNumber,
		CDealNumber:      companyInfo.CDealNumber,
		EReleaseNumber:   companyInfo.EReleaseNumber,
		EDealNumber:      companyInfo.EDealNumber,
		ExpiryDate:       companyInfo.ExpiryDate,
	}, ""
}

type PublicSalesManageListReturn struct {
	List []SalesManageList `json:"SalesManageList"`
}

type SalesManageList struct {
	UserId  int    `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	AddTime string `json:"add_time"`
}

// 公用-销售管理列表
func (this *PublicModel) Public_SalesManageList(userId int) (data *PublicSalesManageListReturn, errMsg string) {
	data = new(PublicSalesManageListReturn)

	// 用户信息
	userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: userId})
	if errMsg != "" {
		return data, errMsg
	}
	if userInfo == nil {
		return data, "该用户不存在"
	}

	// 获取公司员工列表
	sql := `SELECT id, name, email, add_time
			FROM p_user
			WHERE company_id=? AND user_type=0`
	rows, err := db.Db.Query(sql, userInfo.CompanyId)
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
			AddTime: string(value["add_time"]),
		})
	}
	return
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
	// 用户信息
	userInfo, errMsg := userModel.GetUserInfo(&GetUserInfoParameter{UserId: leaderUserId})
	if errMsg != "" {
		return errMsg
	}
	if userInfo == nil {
		return "该用户不存咋"
	}

	// 根据userId是否为0来判断是添加还是编辑
	if userId == 0 { // 添加
		sql := `INSERT INTO p_user(company_id, email, password, name) VALUES(?,?,?,?)`
		_, err := db.Db.Exec(sql, userInfo.CompanyId, email, password, name)
		if err != nil {
			return "添加用户失败"
		}
	} else { // 编辑
		sql := `UPDATE p_user
				SET email=?, password=?, name=?
				WHERE id=?`
		_, err := db.Db.Exec(sql, email, password, name, userId)
		if err != nil {
			return "编辑用户失败"
		}
	}

	return
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
