package v1

import "github.com/gogf/gf/v2/frame/g"

type RegisterDoReq struct {
	// tags openapi
	g.Meta    `path:"/register" method:"post" summary:"用户注册" tags:"注册"`
	Name      string `json:"name#请输入账号" v:"required"`
	Password  string `json:"password#请输入密码" v:"required"`
	Nickename string `json:"nikeName#请输入昵称" v:"required"`
}
type RegisterDoRes struct {
	Result string `json:"result" dc:"注册结果"`
}
