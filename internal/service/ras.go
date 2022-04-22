package service

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"firstproject/internal/model/entity"
	"firstproject/internal/service/internal/dao"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"strings"
)

// RAS管理服务
func Rsa() *sRsa {
	return &sRsa{}
}

type sRsa struct {
}

/**
 * 解密
 */
func (s *sRsa) Decrypt(secretData []byte, name, publicKey string) ([]byte, error) {
	var (
		user *entity.User
		ctx  = gctx.New()
	)
	err := dao.User.Ctx(ctx).Fields("sprivatekey").Where("name=?", name).Limit(1).Scan(&user)
	if err != nil {
		return nil, err
	}
	// 解码将找到下一个PEM格式化块(证书，私钥，etc)
	block, _ := pem.Decode([]byte(user.Sprivatekey))
	var rsaPrivateKey *rsa.PrivateKey
	// 判断私钥类型
	//pkcs1
	if strings.Index(user.Sprivatekey, "BEGIN RSA") > 0 {
		rsaPrivateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else { //pkcs8
		privateKey, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
		rsaPrivateKey = privateKey.(*rsa.PrivateKey)
	}
	// 提取公钥
	block_, _ := pem.Decode([]byte(publicKey))
	publicKey_, _ := x509.ParsePKIXPublicKey(block_.Bytes)
	rsaPublicKey := publicKey_.(*rsa.PublicKey)

	blockLength := rsaPublicKey.N.BitLen() / 8
	// 服务端用给客户端公钥对应的私钥解密（如果公钥被拦截 无法识别用户 需要用到签名）
	if len(secretData) <= blockLength {
		return rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, secretData)
	}

	buffer := bytes.NewBufferString("")

	pages := len(secretData) / blockLength
	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(secretData) {
				continue
			}
			end = len(secretData)
		}

		chunk, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, secretData[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

/**
 * 验签 服务端用客户端事先给服务端的公钥和签名 验签 （验签完成后才会去执行解密数据）
 */
func (s *sRsa) Verify(data []byte, sign []byte, algorithmSign crypto.Hash, name string) bool {
	var (
		user *entity.User
		ctx  = gctx.New()
	)
	err := dao.User.Ctx(ctx).Fields("cpublickey").Where("name=?", name).Limit(1).Scan(&user)
	if err != nil {
		return false
	}
	// 提取公钥
	block, _ := pem.Decode([]byte(user.Cpublickey))
	publicKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	rsaPublicKey := publicKey.(*rsa.PublicKey)
	h := algorithmSign.New()
	h.Write(data)
	// 用c端用户公钥和签名做验签
	// 客户公钥需要保存在服务端
	return rsa.VerifyPKCS1v15(rsaPublicKey, algorithmSign, h.Sum(nil), sign) == nil
}

/**
 * 生成pkcs1格式公钥私钥
 */
func (s *sRsa) CreateKeys(keyLength int) (privateKey, publicKey string) {
	// 随机生成密钥对
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}

	privateKey = string(pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY",
		//pkcs1
		Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
	}))

	derPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return
	}

	publicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}))
	return
}

/**
 * 生成pkcs8格式公钥私钥
 */
func (s *sRsa) CreatePkcs8Keys(keyLength int) (privateKey, publicKey string) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}

	privateKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: s.MarshalPKCS8PrivateKey(rsaPrivateKey),
	}))

	derPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return
	}

	publicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}))
	return
}

func (s *sRsa) MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k
}

// 返回公钥及更新数据
func (s *sRsa) GetKey(ctx context.Context, name string, keyLength int) (publicKey string) {
	privateKey, publicKey := s.CreatePkcs8Keys(keyLength)
	dao.User.Ctx(ctx).Update(fmt.Sprintf("sprivatekey='%s'", privateKey), "name", name)
	return
}

//  push本地公钥
func (s sRsa) Pushkey(ctx context.Context, public, name string) error {
	_, err := dao.User.Ctx(ctx).Update(fmt.Sprintf("cpublickey='%s'", public), "name", name)
	return err
}

// 返回公私钥

func (s *sRsa) GetKeyPri(ctx context.Context, keyLength int) (privateKey, publicKey string) {
	privateKey, publicKey = s.CreatePkcs8Keys(keyLength)
	return
}
