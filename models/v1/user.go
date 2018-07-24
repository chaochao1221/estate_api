package v1

import (
	"estate/db"
	"estate/middleware"
	"estate/pkg/redis"
	"estate/utils"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type UserModel struct{}

type UserLoginReturn struct {
	Authorization string `json:"Authorization"`
}

// 登录
func (this *UserModel) User_Login(email, password string, groupId int) (u *UserLoginReturn, errMsg string) {
	// 根据不同分组验证其相关邮箱密码有效性
	var (
		sql string
		row []map[string][]byte
		err error
	)
	switch groupId {
	case 1: // 本部
		sql = `SELECT id, user_type
			   FROM base_user
			   WHERE email=? AND password=? AND is_del=0`
	case 2: // 中国
		sql = `SELECT id, user_type
			   FROM china_user
			   WHERE email=? AND password=? AND is_del=0`
	case 3: // 日本
		sql = `SELECT u.id, u.user_type, c.expiry_date
			   FROM japan_user u
			   LEFT JOIN japan_company c ON c.id=u.company_id
			   WHERE u.email=? AND u.password=? AND u.is_del=0`
	default:
		return u, "不存在该分组类型"
	}
	row, err = db.Db.Query(sql, email, password)
	if err != nil {
		return u, "获取用户信息失败"
	}
	if len(row) == 0 {
		return u, "邮箱或密码错误"
	}
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

// 生成jwt（24小时有效期）
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

// 登录踢出
func KickOut(k *KickOutParameter) {
	authorization, _ := redis.GetString("GET", "auth#"+k.Email)
	if authorization != "" {
		redis.Do("DEL", authorization)
	}
	redis.Do("SETEX", "auth#"+k.Email, 3600*24, k.JwtString)
	redis.Do("SETEX", k.JwtString, 3600*24, 1)
	return
}

// 更新密码
func (this *UserModel) User_UpdatePassword(email, newPassword string, groupId int) (errMsg string) {
	var tableName string
	switch groupId {
	case 1:
		tableName = `base_user`
	case 2:
		tableName = `china_user`
	case 3:
		tableName = `japan_user`
	default:
		return "不存在该分组类型."
	}
	sql := `UPDATE ` + tableName + ` SET password=? WHERE email=?`
	_, err := db.Db.Exec(sql, newPassword, email)
	if err != nil {
		return "更新密码失败"
	}
	return
}
