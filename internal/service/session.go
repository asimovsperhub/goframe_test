package service

import (
	"context"
	"firstproject/internal/model/entity"
)

type sSession struct {
}

// 目前用jwt-token
// session管理服务
func Session() *sSession {
	// 返回值
	return &sSession{}
}

func (s *sSession) SetUser(ctx context.Context, user *entity.User) error {
	return nil
}
