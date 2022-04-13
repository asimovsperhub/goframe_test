package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// get请求
type LoginIndexReq struct {
	g.Meta `path:"/login" method:"get" summary:"展示登录页面" tags:"登录"`
}
type LoginIndexRes struct {
	//g.Meta `mime:"text/html" type:"string" example:"<html/>"`
	Result string `json:"result" dc:"登录主页"`
}

// post请求
type LoginDoReq struct {
	g.Meta   `path:"/login" method:"post" summary:"执行登录请求" tags:"登录"`
	Name     string `json:"name" v:"required#请输入账号"   dc:"账号"`
	Password string `json:"password" v:"required#请输入密码"   dc:"密码(明文)"`
}
type LoginDoRes struct {
	Result string `json:"result" dc:"登录结果"`
}
