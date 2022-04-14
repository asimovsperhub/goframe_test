package model

// UserLoginInput 用户登录
type UserLoginInput struct {
	Name     string // 账号
	Password string // 密码(明文)
}

// UserRegisterInput 用户注册
type UserRegisterInput struct {
	Name     string // 账号
	NikeName string //昵称
	Password string // 密码(明文)
}
