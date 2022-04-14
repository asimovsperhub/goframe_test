package v1

import "github.com/gogf/gf/v2/frame/g"

// get请求
type UserInfoReq struct {
	g.Meta `path:"/user/info" method:"get" summary:"用户信息" tags:"用户信息"`
}
type UserInfoRes struct {
	//g.Meta `mime:"text/html" type:"string" example:"<html/>"`
	Result string `json:"result" dc:"用户信息"`
}
