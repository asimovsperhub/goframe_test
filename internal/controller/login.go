package controller

import (
	"context"
	v1 "firstproject/api/v1"
	"firstproject/internal/model"
	"firstproject/internal/service"
)

// 登录包管理
// 路由绑定用到 对外调用
var (
	Login = cLogin{}
)

type cLogin struct{}

//  登录首页
func (c *cLogin) Index(ctx context.Context, req *v1.LoginIndexReq) (res *v1.LoginIndexRes, err error) {
	res = &v1.LoginIndexRes{}
	res.Result = "login index"
	return
}

// 登录处理
func (c *cLogin) Login(ctx context.Context, req *v1.LoginDoReq) (res *v1.LoginDoRes, err error) {
	if err = service.User().Login(ctx, model.UserLoginInput{Name: req.Name, Password: req.Password}); err != nil {
		res = &v1.LoginDoRes{}
		res.Result = "login failed"
	} else {
		res = &v1.LoginDoRes{}
		res.Token, res.Expire = service.Auth().LoginHandler(ctx)
		res.Result = "login success"
		return
	}
	return
}
