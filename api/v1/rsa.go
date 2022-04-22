package v1

// get请求
import "github.com/gogf/gf/v2/frame/g"

type GetKeyReq struct {
	g.Meta `path:"/get_key" method:"get" summary:"获取公钥" tags:"RSA"`
	Name   string `json:"name" v:"required#账号"   dc:"账号"`
}
type GetKeyRes struct {
	//g.Meta `mime:"text/html" type:"string" example:"<html/>"`
	Result string `json:"result" dc:"返回公钥"`
	Public string `json:"public"  dc:"公钥"`
}

type GetKeyPriReq struct {
	g.Meta `path:"/get_key_pri" method:"get" summary:"获取公钥" tags:"RSA"`
}
type GetKeyPriRes struct {
	//g.Meta `mime:"text/html" type:"string" example:"<html/>"`
	Result  string `json:"result" dc:"返回公钥"`
	Public  string `json:"public"  dc:"公钥"`
	Private string `json:"private"  dc:"公钥"`
}

type PostPublicKeyReq struct {
	g.Meta `path:"/push_key" method:"post" summary:"用户上传公钥" tags:"RSA"`
	Public string `json:"public" v:"required#公钥"   dc:"公钥"`
	Name   string `json:"name" v:"required#账号"   dc:"账号"`
}
type PostPublicKeyRes struct {
	Result string `json:"result" dc:"结果"`
}

type VerifyReq struct {
	g.Meta `path:"/verify" method:"post" summary:"验签" tags:"RSA"`
	Sign   string `json:"sign" v:"required#签名"   dc:"签名"`
	Name   string `json:"name" v:"required#账号"   dc:"账号"`
	Data   string `json:"data" v:"required#签署信息"   dc:"签署信息"`
}
type VerifyRes struct {
	Result bool `json:"result" dc:"验签结果"`
}
type DecryptReq struct {
	g.Meta    `path:"/decrypt" method:"post" summary:"解密" tags:"RSA"`
	Secret    string `json:"secret" v:"required#密文"   dc:"密文"`
	Name      string `json:"name" v:"required#账号"   dc:"账号"`
	PublicKey string `json:"publickey" v:"required#公钥"   dc:"公钥"`
}
type DecryptRes struct {
	Result string `json:"result" dc:"解密结果"`
}
