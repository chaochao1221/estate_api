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
