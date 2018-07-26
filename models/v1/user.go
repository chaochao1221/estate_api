package v1

import (
	"estate/db"
	"estate/middleware"
	"estate/pkg/redis"
	"estate/utils"
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type UserModel struct{}

type UserLoginReturn struct {
	Authorization string `json:"Authorization"`
}

// 用户-登录
func (this *UserModel) User_Login(email, password string) (u *UserLoginReturn, errMsg string) {
	// 根据不同分组验证其相关邮箱密码有效性
	sql := `SELECT u.id, u.group_id, u.user_type, c.expiry_date
			FROM p_user u
			LEFT JOIN p_company c ON c.id=u.company_id
			WHERE u.email=? AND u.password=?`
	row, err := db.Db.Query(sql, email, password)
	if err != nil {
		return u, "获取用户信息失败"
	}
	if len(row) == 0 {
		return u, "邮箱或密码错误"
	}
	groupId, _ := strconv.Atoi(string(row[0]["group_id"]))
	if groupId == 3 {
		expiryDate := string(row[0]["expiry_date"])
		if time.Now().Format("2006-01-02 15:04:05") > expiryDate {
			return u, "帐号已过期"
		}
	}

	// 生成jwt
	jwtString := CreateJwt(&CreateJwtParameter{
		UserId:   utils.Str2int(string(row[0]["id"])),
		UserType: utils.Str2int(string(row[0]["user_type"])),
		GroupId:  groupId,
	})

	// 登录踢出
	KickOut(&KickOutParameter{
		JwtString: jwtString,
		Email:     email,
	})

	// 返回数据
	return &UserLoginReturn{
		Authorization: "Bearer " + jwtString,
	}, ""
}

type CreateJwtParameter struct {
	UserId   int
	UserType int
	GroupId  int
}

/*
* @Title CreateJwt
* @Description 生成jwt（有效期：24小时）
* @Parameter c *CreateJwtParameter
* @Return jwtString string
 */
func CreateJwt(c *CreateJwtParameter) (jwtString string) {
	claims := middleware.Claims{
		"estate",
		c.UserId,
		c.UserType,
		c.GroupId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(24) * time.Hour).Unix(),
		},
	}
	jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, _ = jwttoken.SignedString([]byte(middleware.AuthKey))
	return
}

type KickOutParameter struct {
	JwtString string
	Email     string
}

/*
* @Title KickOut
* @Description 登录踢出
* @Parameter k *KickOutParameter
 */
func KickOut(k *KickOutParameter) {
	authorization, _ := redis.GetString("GET", "auth#"+k.Email)
	if authorization != "" {
		redis.Do("DEL", authorization)
	}
	redis.Do("SETEX", "auth#"+k.Email, 3600*24, k.JwtString)
	redis.Do("SETEX", k.JwtString, 3600*24, 1)
	return
}

/*
* @Title User_UpdatePassword
* @Description 更新用户密码
* @Parameter email string
* @Parameter newPassword string
* @Return errMsg string
 */
func (this *UserModel) User_UpdatePassword(email, newPassword string) (errMsg string) {
	sql := `UPDATE p_user SET password=? WHERE email=?`
	_, err := db.Db.Exec(sql, newPassword, email)
	if err != nil {
		return "更新密码失败"
	}
	return
}

type UserInfoReturn struct {
	GroupId    int    `json:"group_id"`
	UserId     int    `json:"user_id"`
	UserType   int    `json:"user_type"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	IsNotified int    `json:"is_notified"`
}

// 用户-信息
func (this *UserModel) User_Info(userId, groupId int) (u *UserInfoReturn, errMsg string) {
	// 获取用户信息
	userInfo, errMsg := this.GetUserInfo(&GetUserInfoParameter{UserId: userId})
	if errMsg != "" {
		return u, errMsg
	}
	if userInfo == nil {
		return u, "该用户不存在"
	}

	// 获取本部基础信息
	baseInfo, errMsg := this.GetBaseInfo()
	if errMsg != "" {
		return u, errMsg
	}

	return &UserInfoReturn{
		GroupId:    groupId,
		UserId:     userId,
		UserType:   userInfo.UserType,
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		IsNotified: baseInfo.IsNotified,
	}, ""
}

type GetUserInfoParameter struct {
	UserId int
	Email  string
}

type GetUserInfoReturn struct {
	CompanyId int
	UserId    int
	UserType  int
	Email     string
	Password  string
	Name      string
	Telephone string
	Fax       string
	AddTime   string
}

/*
* @Title GetUserInfo
* @Description 获取用户信息（根据调用者传参是用户id还是邮箱来查询用户信息）
* @Parameter userInfo *GetUserInfoParameter
* @Return u *GetUserInfoReturn
* @Return errMsg string
 */
func (this *UserModel) GetUserInfo(userInfo *GetUserInfoParameter) (u *GetUserInfoReturn, errMsg string) {
	// 判断是根据用户id或者邮箱查询用户信息
	var where string
	if userInfo.UserId > 0 {
		where = `id=` + strconv.Itoa(userInfo.UserId)
	} else {
		where = `email='` + userInfo.Email + `'`
	}

	// 查询用户信息
	sql := `SELECT id, company_id, user_type, email, password, name, telephone, fax, add_time
			FROM p_user
			WHERE ` + where
	row, err := db.Db.Query(sql)
	if err != nil {
		return u, "获取用户信息失败"
	}
	if len(row) == 0 {
		return u, "该帐号不存在"
	}

	// 返回数据
	return &GetUserInfoReturn{
		CompanyId: utils.Str2int(string(row[0]["company_id"])),
		UserId:    utils.Str2int(string(row[0]["id"])),
		UserType:  utils.Str2int(string(row[0]["user_type"])),
		Email:     string(row[0]["email"]),
		Password:  string(row[0]["password"]),
		Name:      string(row[0]["name"]),
		Telephone: string(row[0]["telephone"]),
		Fax:       string(row[0]["fax"]),
		AddTime:   string(row[0]["add_time"]),
	}, ""
}

type GetCompanyInfoReturn struct {
	GroupId          int
	RegionId         int
	Name             string
	Adress           string
	CRecommendNumber int
	CDealNumber      int
	EReleaseNumber   int
	EDealNumber      int
	ExpiryDate       string
}

/*
* @Title GetCompanyInfo
* @Description 获取公司信息
* @Parameter companyId int
* @Return data *GetCompanyInfoReturn
* @Return errMsg string
 */
func (this *UserModel) GetCompanyInfo(companyId int) (data *GetCompanyInfoReturn, errMsg string) {
	sql := `SELECT group_id, region_id, name, adress, c_recommend_number, c_deal_number, e_release_number, e_deal_number, expiry_date
			FROM p_company
			WHERE id=?`
	row, err := db.Db.Query(sql, companyId)
	if err != nil {
		return data, "获取公司信息失败"
	}
	return &GetCompanyInfoReturn{
		GroupId:          utils.Str2int(string(row[0]["group_id"])),
		RegionId:         utils.Str2int(string(row[0]["region_id"])),
		Name:             string(row[0]["name"]),
		Adress:           string(row[0]["adress"]),
		CRecommendNumber: utils.Str2int(string(row[0]["c_recommend_number"])),
		CDealNumber:      utils.Str2int(string(row[0]["c_deal_number"])),
		EReleaseNumber:   utils.Str2int(string(row[0]["e_release_number"])),
		EDealNumber:      utils.Str2int(string(row[0]["e_deal_number"])),
		ExpiryDate:       string(row[0]["expiry_date"]),
	}, ""
}

type GetBaseInfoReturn struct {
	ServiceFee       int
	FixedFee         int
	ExciseFee        string
	ProtectionPeriod int
	IsNotified       int
}

/*
* @Title GetBaseInfo
* @Description 获取本部基础信息
* @Return data *GetBaseInfoReturn
* @Return errMsg string
 */
func (this *UserModel) GetBaseInfo() (data *GetBaseInfoReturn, errMsg string) {
	sql := `SELECT service_fee, fixed_fee, excise_fee, protection_period, is_notified
			FROM base_info
			WHERE id=1`
	row, err := db.Db.Query(sql)
	if err != nil {
		return data, "获取本部基础信息失败"
	}
	return &GetBaseInfoReturn{
		ServiceFee:       utils.Str2int(string(row[0]["service_fee"])),
		FixedFee:         utils.Str2int(string(row[0]["fixed_fee"])),
		ExciseFee:        string(row[0]["excise_fee"]),
		ProtectionPeriod: utils.Str2int(string(row[0]["protection_period"])),
		IsNotified:       utils.Str2int(string(row[0]["is_notified"])),
	}, ""
}

// 用户-修改密码
func (this *UserModel) User_ModifyPassword(userId int, oldPassword, newPassword string) (errMsg string) {
	// 获取用户信息
	userInfo, errMsg := this.GetUserInfo(&GetUserInfoParameter{UserId: userId})
	if errMsg != "" {
		return errMsg
	}
	if userInfo == nil {
		return "该用户不存在"
	}

	// 验证原密码是否正确
	if oldPassword != userInfo.Password {
		return "原密码错误"
	}

	// 更新密码
	errMsg = this.User_UpdatePassword(userInfo.Email, newPassword)
	if errMsg != "" {
		return
	}
	return
}

// 用户-重置密码
func (this *UserModel) User_ResetPassword(email string) (errMsg string) {
	// 获取用户信息
	userInfo, errMsg := this.GetUserInfo(&GetUserInfoParameter{Email: email})
	if errMsg != "" {
		return errMsg
	}
	if userInfo == nil {
		return "该帐号不存在"
	}

	// 获取公司信息
	companyInfo, errMsg := this.GetCompanyInfo(userInfo.CompanyId)
	if errMsg != "" {
		return errMsg
	}

	if companyInfo.GroupId == 3 && companyInfo.ExpiryDate < time.Now().Format("2006-01-02 15:04:05") {
		return "该帐号已过期"
	}

	// 在redis里重置密码，有效期24小时
	newPassword := string(utils.Krand(6, 0))
	_, err := redis.Do("SETEX", "resetPassword#"+email, 60*60*24, newPassword)
	if err != nil {
		return "重置密码失败"
	}

	// 向邮箱发送重置后的密码
	fmt.Println(newPassword)

	return
}
