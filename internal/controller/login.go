package controller

import (
	"context"
	v1 "firstproject/api/v1"
	"firstproject/internal/model"
	"firstproject/internal/service"
	"fmt"
)

// 登录包管理
// 路由绑定用到 对外调用
var (
	Login = cLogin{}
)

type cLogin struct{}

//  登录首页
func (c *cLogin) Index(ctx context.Context, req *v1.LoginIndexReq) (res *v1.LoginIndexRes, err error) {
	return &v1.LoginIndexRes{"login index"}, nil
}

// 登录处理
func (c *cLogin) Login(ctx context.Context, req *v1.LoginDoReq) (res *v1.LoginDoRes, err error) {
	if err = service.User().Login(ctx, model.UserLoginInput{Name: req.Name, Password: req.Password}); err != nil {
		return &v1.LoginDoRes{fmt.Sprintf("login failed:%s", err)}, nil
	} else {
		return &v1.LoginDoRes{"login success"}, nil
	}
}
