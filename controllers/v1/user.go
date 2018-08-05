package v1

import (
	"estate/middleware"
	"estate/models/v1"
	"estate/pkg/redis"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var userModel = new(v1.UserModel)

// 用户-路由
func User(parentRoute *gin.RouterGroup) {
	router := parentRoute.Group("/user")
	router.POST("/login", User_Login)                  // 2.1 用户-登录
	router.POST("/reset_password", User_ResetPassword) // 2.4 用户-重置密码
	router.Use(middleware.Auth())
	router.GET("/info", User_Info)                       // 2.2 用户-信息
	router.POST("/modify_password", User_ModifyPassword) // 2.3 用户-修改密码
	router.DELETE("/logout", User_Logout)                // 2.5 用户-注销
}

// 用户-登录
func User_Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
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
	if newPassword == password {
		// 重置密码
		errMsg := userModel.User_UpdatePassword(email, newPassword)
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
	data, errMsg := userModel.User_Login(email, password)
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

// 用户-信息
func User_Info(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	groupId, _ := strconv.Atoi(c.Request.Header.Get("group_id"))
	if userId == 0 || groupId == 0 {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 获取用户信息
	data, errMsg := userModel.User_Info(userId, groupId)
	if errMsg != "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  errMsg,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
	return
}

// 用户-修改密码
func User_ModifyPassword(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("user_id"))
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	if userId == 0 || oldPassword == "" || newPassword == "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 修改密码
	errMsg := userModel.User_ModifyPassword(userId, oldPassword, newPassword)
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
	})
}

// 用户-重置密码
func User_ResetPassword(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "参数错误",
		})
		return
	}

	// 重置密码
	errMsg := userModel.User_ResetPassword(email)
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
	})
	return
}

// 用户-注销
func User_Logout(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	r, _ := regexp.Compile("^Bearer (.+)$")
	match := r.FindStringSubmatch(authHeader)
	tokenString := match[1]
	_, err := redis.Do("DEL", tokenString)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 1010,
			"msg":  "注销失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
	})
	return
}
