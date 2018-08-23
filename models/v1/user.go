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
func (this *UserModel) User_Login(email, password string) (data *UserLoginReturn, errMsg string) {
	// 用户信息
	userInfo, errMsg := this.GetUserInfo(&GetUserInfoParameter{Email: email})
	if errMsg != "" {
		return data, errMsg
	}

	// 公司信息
	companyInfo, errMsg := this.GetCompanyInfo(userInfo.CompanyId)
	if errMsg != "" {
		return data, errMsg
	}

	// 判断帐号是否过期
	if companyInfo.GroupId == 3 && time.Now().Format("2006-01-02") > companyInfo.ExpiryDate {
		return data, "The account has expired"
	}

	// 验证密码
	if !utils.CheckPassword(userInfo.Password, password) {
		return data, "Password error"
	}

	// 生成jwt
	jwtString := CreateJwt(&CreateJwtParameter{
		UserId:   userInfo.UserId,
		UserType: userInfo.UserType,
		GroupId:  companyInfo.GroupId,
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
	_, err := db.Db.Exec(sql, string(utils.HashPassword(newPassword)), email)
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

	return &UserInfoReturn{
		GroupId:    groupId,
		UserId:     userId,
		UserType:   userInfo.UserType,
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		IsNotified: userInfo.IsNotified,
	}, ""
}

type GetUserInfoParameter struct {
	UserId int
	Email  string
}

type GetUserInfoReturn struct {
	CompanyId  int
	UserId     int
	UserType   int
	Email      string
	Password   string
	Name       string
	Telephone  string
	Fax        string
	IsNotified int
	AddTime    string
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
	sql := `SELECT id, company_id, user_type, email, password, name, telephone, fax, is_notified, add_time
			FROM p_user
			WHERE ` + where
	row, err := db.Db.Query(sql)
	if err != nil {
		return u, "获取用户信息失败"
	}
	if len(row) == 0 {
		return u, "The account does not exist"
	}

	// 返回数据
	return &GetUserInfoReturn{
		CompanyId:  utils.Str2int(string(row[0]["company_id"])),
		UserId:     utils.Str2int(string(row[0]["id"])),
		UserType:   utils.Str2int(string(row[0]["user_type"])),
		Email:      string(row[0]["email"]),
		Password:   string(row[0]["password"]),
		Name:       string(row[0]["name"]),
		Telephone:  string(row[0]["telephone"]),
		Fax:        string(row[0]["fax"]),
		IsNotified: utils.Str2int(string(row[0]["is_notified"])),
		AddTime:    string(row[0]["add_time"]),
	}, ""
}

type GetCompanyInfoReturn struct {
	GroupId         int
	RegionId        int
	Name            string
	Adress          string
	RecommendNumber int // 推荐客户数
	ReleaseNumber   int // 发布房源数
	ButtNumber      int // 对接客户数
	DealNumber      int // 成交客户/房源数
	ExpiryDate      string
}

/*
* @Title GetCompanyInfo
* @Description 获取公司信息
* @Parameter companyId int
* @Return data *GetCompanyInfoReturn
* @Return errMsg string
 */
func (this *UserModel) GetCompanyInfo(companyId int) (data *GetCompanyInfoReturn, errMsg string) {
	sql := `SELECT group_id, region_id, name, adress, recommend_number, release_number, butt_number, deal_number, expiry_date
			FROM p_company
			WHERE id=?`
	row, err := db.Db.Query(sql, companyId)
	if err != nil {
		return data, "获取公司信息失败"
	}
	return &GetCompanyInfoReturn{
		GroupId:         utils.Str2int(string(row[0]["group_id"])),
		RegionId:        utils.Str2int(string(row[0]["region_id"])),
		Name:            string(row[0]["name"]),
		Adress:          string(row[0]["adress"]),
		RecommendNumber: utils.Str2int(string(row[0]["recommend_number"])),
		ReleaseNumber:   utils.Str2int(string(row[0]["release_number"])),
		ButtNumber:      utils.Str2int(string(row[0]["butt_number"])),
		DealNumber:      utils.Str2int(string(row[0]["deal_number"])),
		ExpiryDate:      string(row[0]["expiry_date"]),
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
	sql := `SELECT service_fee, fixed_fee, excise_fee, protection_period
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
	}, ""
}

// 用户-修改密码
func (this *UserModel) User_ModifyPassword(userId int, oldPassword, newPassword string) (errMsg string) {
	// 获取用户信息
	userInfo, errMsg := this.GetUserInfo(&GetUserInfoParameter{UserId: userId})
	if errMsg != "" {
		return errMsg
	}

	// 验证原密码是否正确
	if !utils.CheckPassword(userInfo.Password, oldPassword) {
		return "Original password error"
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

	// 获取公司信息
	companyInfo, errMsg := this.GetCompanyInfo(userInfo.CompanyId)
	if errMsg != "" {
		return errMsg
	}

	// 判断帐号是否过期
	if companyInfo.GroupId == 3 && companyInfo.ExpiryDate < time.Now().Format("2006-01-02") {
		return "The account has expired"
	}

	// 在redis里重置密码，有效期24小时
	newPassword := string(utils.Krand(6, 0))
	_, err := redis.Do("SETEX", "resetPassword#"+email, 60*60*24, newPassword)
	if err != nil {
		return "重置密码失败"
	}

	// 向邮箱发送重置后的密码
	fmt.Println(newPassword)
	if errMsg = utils.SendEmail(email, newPassword, "", 2); errMsg != "" {
		return errMsg
	}

	return
}
