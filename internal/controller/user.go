package controller

import (
	"context"
	v1 "firstproject/api/v1"
	"fmt"
)

var (
	User = cUser{}
)

type cUser struct {
}

func (*cUser) Info(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error) {
	return &v1.UserInfoRes{fmt.Sprintf("user info")}, nil
}
