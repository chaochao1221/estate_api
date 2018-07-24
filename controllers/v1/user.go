package v1

import (
	"estate/middleware"
	"estate/models/v1"
	"estate/pkg/redis"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

var userModel = new(v1.UserModel)

// 路由
func User(parentRoute *gin.RouterGroup) {
	router := parentRoute.Group("/user")
	router.POST("login", User_Login)
	router.Use(middleware.Auth())
}

// 登录
func User_Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	groupId, _ := strconv.Atoi(c.PostForm("group_id"))
	if email == "" || password == "" || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 先查看redis里是否设有需要重置的密码（resetPassword#email=>password）
	// 若存在重置密码且与传入的密码一致，则修改库中密码，删除redis重置密码，并验证库中邮箱密码的有效性
	// 若重置密码与传入的密码不一致，则直接去数据库中验证邮箱密码的有效性
	newPassword, _ := redis.GetString("GET", "resetPassword#"+email)
	fmt.Println("newPassword", newPassword)
	if newPassword == password {
		// 重置密码
		errMsg := userModel.User_UpdatePassword(email, newPassword, groupId)
		if errMsg != "" {
			c.JSON(400, gin.H{
				"code": 0,
				"msg":  errMsg,
			})
			return
		}

		// 清空重置密码缓存
		_, err := redis.Do("DEL", "resetPassword#"+email)
		if err != nil {
			c.JSON(400, gin.H{
				"code": 0,
				"msg":  err.Error(),
			})
		}
	}
	// 验证邮箱密码有效性
	data, errMsg := userModel.User_Login(email, password, groupId)
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
		"data": data,
	})
	return
}
