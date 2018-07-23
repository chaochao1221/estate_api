//error model

package models

// api 错误返回
type APIError struct {
	Code uint   `json:"code" description:"错误编码"`
	Msg  string `json:"msg"  description:"错误描述"`
}
