package controller

import (
	"context"
	"crypto"
	v1 "firstproject/api/v1"
	"firstproject/internal/service"
)

var (
	Rsa = cRsa{}
)

type cRsa struct{}

// 获取公钥
func (c *cRsa) GetKey(ctx context.Context, req *v1.GetKeyReq) (res *v1.GetKeyRes, err error) {
	//context.WithValue(ctx, "name", req.Name)
	res = &v1.GetKeyRes{}
	pub := service.Rsa().GetKey(ctx, req.Name, 2048)
	res.Public = pub
	res.Result = "success"
	return
}

// 生成公私钥
func (c *cRsa) GetKeyPri(ctx context.Context, req *v1.GetKeyPriReq) (res *v1.GetKeyPriRes, err error) {
	//context.WithValue(ctx, "name", req.Name)
	res = &v1.GetKeyPriRes{}
	pri, pub := service.Rsa().GetKeyPri(ctx, 2048)
	res.Public = pub
	res.Private = pri
	res.Result = "success"
	return
}

// push 私钥
func (c *cRsa) PushKey(ctx context.Context, req *v1.PostPublicKeyReq) (res *v1.PostPublicKeyRes, err error) {
	res = &v1.PostPublicKeyRes{}
	err = service.Rsa().Pushkey(ctx, req.Public, req.Name)
	if err != nil {
		res.Result = "push public failed"
		return
	} else {
		res.Result = "push public success"
		return
	}

}

// 解密
func (c *cRsa) Decrypt(ctx context.Context, req *v1.DecryptReq) (res *v1.DecryptRes, err error) {
	res = &v1.DecryptRes{}
	secret := []byte(req.Secret)
	if secret_, err := service.Rsa().Decrypt(secret, req.Name, req.PublicKey); err != nil {
		res.Result = "Decrypt failed"
		return res, err
	} else {
		res.Result = "Decrypt success:" + string(secret_)
		return res, nil
	}
}

// 验签
func (c *cRsa) Verify(ctx context.Context, req *v1.VerifyReq) (res *v1.VerifyRes, err error) {
	res = &v1.VerifyRes{}
	sign := []byte(req.Sign)
	data := []byte(req.Data)
	result := service.Rsa().Verify(data, sign, crypto.SHA256, req.Name)
	res.Result = result
	return
}
