package service

import (
	"context"
	"firstproject/internal/model"
	"firstproject/internal/model/entity"
	"firstproject/internal/service/internal/dao"
	"fmt"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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
	//  去数据库查询用户信息
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
	// set token
	if err := Session().SetUser(ctx, userEntity); err != nil {
		return err
	}
	return nil
}

// 用户注册。
// 事务使用函数f包装事务逻辑。 失败回滚
func (s *sUser) Register(ctx context.Context, in model.UserRegisterInput) error {
	return dao.User.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		var user *entity.User
		// Struct将params键值对映射到相应的Struct对象的属性
		// 将传进来的值 映射到user对象
		if err := gconv.Struct(in, &user); err != nil {
			return err
		}
		if err := s.CheckPassportUnique(ctx, user.Name); err != nil {
			return err
		}
		if err := s.CheckNicknameUnique(ctx, user.Nickename); err != nil {
			return err
		}
		user.Password = s.EncryptPassword(user.Name, user.Password)
		//  写入数据库
		fmt.Println(*user)
		_, err := dao.User.Ctx(ctx).Data(user).OmitEmpty().Save()
		return err
	})
}

// 检测给定的账号是否唯一
func (s *sUser) CheckPassportUnique(ctx context.Context, name string) error {
	n, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Name, name).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return gerror.Newf(`账号"%s"已被占用`, name)
	}
	return nil
}

// 检测给定的昵称是否唯一
func (s *sUser) CheckNicknameUnique(ctx context.Context, nickname string) error {
	// 查询数据库改昵称数量
	n, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Nickename, nickname).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return gerror.Newf(`昵称"%s"已被占用`, nickname)
	}
	return nil
}

func (s *sUser) GetUserInfo(ctx context.Context, in model.UserLoginInput) map[string]interface{} {
	return g.Map{
		"id":   1,
		"name": in.Name,
	}
}
