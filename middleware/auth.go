package middleware

import (
	"estate/pkg/redis"
	"log"
	"regexp"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const AuthKey = "123456"

type Claims struct {
	Application string `json:"application"`
	UserType    string `json:"user_type"`
	UserId      int    `json:"user_id"`
	GroupId     int    `json:"group_id"`
	jwt.StandardClaims
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		r, _ := regexp.Compile("^Bearer (.+)$")

		match := r.FindStringSubmatch(authHeader)
		if len(match) == 0 {
			c.JSON(401, gin.H{
				"code": 1010,
				"msg":  "token is null!",
			})
			c.Abort()
			return
		}

		tokenString := match[1]
		t, err := redis.GetString("GET", tokenString)
		if t == "" || err != nil {
			c.JSON(401, gin.H{
				"code": 1010,
				"msg":  "token is expired",
			})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(AuthKey), nil
		})
		if err != nil {
			log.Printf("error: %s", err.Error())
			c.JSON(401, gin.H{
				"code": 1010,
				"msg":  "token is error",
			})
			c.Abort()
			return
		} else {
			if claims, ok := token.Claims.(*Claims); ok && token.Valid {
				c.Request.Header.Set("user_id", strconv.Itoa(claims.UserId))
				c.Request.Header.Set("user_type", claims.UserType)
				c.Request.Header.Set("group_id", strconv.Itoa(claims.GroupId))
			} else {
				c.JSON(401, gin.H{
					"code": 1010,
					"msg":  "token is error!",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
