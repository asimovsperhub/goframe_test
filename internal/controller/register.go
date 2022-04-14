package controller

import (
	"context"
	v1 "firstproject/api/v1"
	"firstproject/internal/model"
	"firstproject/internal/service"
	"fmt"
)

// 注册包管理
// 路由绑定用到 对外调用
var (
	Register = cRegister{}
)

type cRegister struct {
}

func (*cRegister) Register(ctx context.Context, req *v1.RegisterDoReq) (res *v1.RegisterDoRes, err error) {
	// 注册 service
	if err = service.User().Register(ctx, model.UserRegisterInput{
		Name:     req.Name,
		Password: req.Password,
		NikeName: req.NikeName,
	}); err != nil {
		return &v1.RegisterDoRes{fmt.Sprintf("register failed:%s", err)}, nil
	} else {
		// 自动登录
		err = service.User().Login(ctx, model.UserLoginInput{
			Name:     req.Name,
			Password: req.Password,
		})
		if err != nil {
			return &v1.RegisterDoRes{fmt.Sprintf("register success login failed %s", err)}, nil
		} else {
			return &v1.RegisterDoRes{fmt.Sprintf("register success login success")}, nil
		}
	}
}
