package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const APP_VER = "1.0"

func Default(c *gin.Context) {
	sr := time.Now().Unix()
	message := `You Win!
ServerName: gin
estate房产管理 API Version: ` + APP_VER + `
RunMode: ` + gin.Mode() + `
ServerTime: ` + time.Unix(sr, 0).Format("2006-01-02 15:04:05") + `(` + strconv.Itoa(int(sr)) + `)`
	c.String(http.StatusOK, message)
}
