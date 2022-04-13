package service

import (
	"context"
	"firstproject/internal/model"
	"firstproject/internal/model/entity"
	"firstproject/internal/service/internal/dao"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sUser struct{}

// 用户管理服务
func User() *sUser {
	return &sUser{}
}

// 根据账号和密码查询用户信息，一般用于账号密码登录。
// 注意password参数传入的是按照相同加密算法加密过后的密码字符串。
// 返回*entity.User 数据字段 查询结果
func (s *sUser) GetUserByPassportAndPassword(ctx context.Context, name, password string) (user *entity.User, err error) {
	err = dao.User.Ctx(ctx).Where(g.Map{
		dao.User.Columns().Name:     name,
		dao.User.Columns().Password: password,
	}).Scan(&user)
	return
}

// 将密码按照内部算法进行加密
func (s *sUser) EncryptPassword(name, password string) string {
	return gmd5.MustEncrypt(name + password)
}

// 执行登录
func (s *sUser) Login(ctx context.Context, in model.UserLoginInput) error {
	//  去数据库查询用户名密码
	userEntity, err := s.GetUserByPassportAndPassword(
		ctx,
		in.Name,
		// 加密后的密码
		s.EncryptPassword(in.Name, in.Password),
	)
	if err != nil {
		return err
	}
	if userEntity == nil {
		// 给定err
		return gerror.New(`账号或密码错误`)
	}
	return nil
}
