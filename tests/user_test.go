// package tests

// import (
// 	"estate/controllers/v1"
// 	"estate/db"
// 	"estate/pkg"
// 	"fmt"
// 	"net/http/httptest"
// 	"net/url"
// 	"strings"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// func SetupRouter() *gin.Engine {
// 	router := gin.Default()
// 	gin.SetMode(gin.TestMode)
// 	api_v1 := router.Group("/v1")
// 	{
// 		v1.User(api_v1)
// 	}
// 	return router
// }

// func main() {
// 	r := SetupRouter()
// 	r.Run()
// }

// func TestIntPKG(t *testing.T) {
// 	pkg.Init()
// }

// func TestIntDB(t *testing.T) {
// 	db.Init()
// }

// func TestVcodeTest(t *testing.T) {
// 	testRouter := SetupRouter()
// 	params := url.Values{}
// 	params.Add("mobile", "13921275871")
// 	r := httptest.NewRequest("POST", "/v1/user/vcode_test", strings.NewReader(params.Encode()))
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	w := httptest.NewRecorder()
// 	testRouter.ServeHTTP(w, r)
// 	fmt.Println(w.Body)
// 	assert.Equal(t, w.Code, 201)
// }
