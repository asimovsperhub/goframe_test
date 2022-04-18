package controller

import (
	"context"
	v1 "firstproject/api/v1"
	"firstproject/internal/service"
	"fmt"
)

var (
	Geth = cGeth{}
)

type cGeth struct {
}

func (*cGeth) Conn(ctx context.Context, req *v1.GethReq) (res *v1.GethRes, err error) {
	if result, err := service.Gethser().DoFunc(ctx); err != nil {
		return &v1.GethRes{fmt.Sprintf("%s", err)}, nil
	} else {
		return &v1.GethRes{fmt.Sprintf(result)}, nil
	}
}
