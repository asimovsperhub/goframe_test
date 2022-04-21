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

// jwt用户私有信息
func (s *sUser) GetUserInfo(ctx context.Context, in model.UserLoginInput) map[string]interface{} {
	return g.Map{
		"id":   1,
		"name": in.Name,
	}
}
func (s *sUser) UpdateUser(ctx context.Context) {
	// 直接update（要改的数据 ,条件字段,value）
	g.Model("user").Update("name='test'", "name", "asimov123")
	// 要改的数据 where 过滤条件
	g.Model("user").Data("name='test'").Where("name", "asimov123").Update()
	// Counter

}

func (s *sUser) DeleteUser(ctx context.Context) {
	dao.User.Ctx(ctx).Where("name", "asimov123").Delete()
	//
	g.Model("user").Where("name", "asimov123").Delete()
	// DELETE FROM `user` ORDER BY `login_time` asc LIMIT 10
	g.Model("user").Order("id asc").Limit(10).Delete()
	// 直接调delete(条件) DELETE FROM `user` WHERE `score`<60
	g.Model("user").Delete("status =", 0)

}

func (s *sUser) FindUser(ctx context.Context) (user []*entity.User, err error) {
	// select  *  from user where status=1 limit 0,1(其实索引,限制数量)
	g.Model("user").Where("status=?", 1).Limit(0, 1).All()
	//
	g.Model("user").Where("status=1").Limit(0, 1).All()
	// 指定字段
	g.Model("user").Fields("name,status").Where("status=?", 1).Limit(0, 1).All()
	// and
	g.Model("user").Where("name", "asimov123").Where("status", 1).One()
	// where + slice  变量查询
	g.Model("user").Where("status=? AND name like ?", g.Slice{1, "asimov%"}).All()
	// where + map
	g.Model("user").Where(g.Map{"status": 1, "name like": "asimov"}).All()
	// wheref
	g.Model("user").Wheref("status=? and name like in (?)", 1, g.Slice{"asimov", "asimov123"})
	// wherepri 主键查询,对该表的主键智能识别 id
	g.Model("user").WherePri(g.Slice{1, 2, 3})
	// 嵌套查询
	g.Model("user").Wheref("status in (?)", g.Model("user").Fields("status").Wheref("name in (?)", g.Slice{"asimov", "asimov123"})).Scan(&user)
	// distinct 查询
	// 定义返回结果
	//var(
	//	user []*entity.User
	//)
	dao.User.Ctx(ctx).Wheref("status=? and name like in (?)", 1, g.Slice{"asimov", "asimov123"}).Scan(&user)
	res, err := dao.User.Ctx(ctx).Wheref("status=? and name like in (?)", 1, g.Slice{"asimov", "asimov123"}).All()
	//
	if len(res) == 0 || res.IsEmpty() {
		return
	} else {
		println(user)
	}

	return
}
