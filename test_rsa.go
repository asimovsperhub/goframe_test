package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"
)

// RAS管理服务
func Rsa(publicKey, privateKey string) *sRsa {
	rsaObj := &sRsa{
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	rsaObj.init()

	return rsaObj
}

type sRsa struct {
	privateKey string
	publicKey  string
	// 私钥
	rsaPrivateKey *rsa.PrivateKey
	// 公钥
	rsaPublicKey *rsa.PublicKey
}

// 初始化
func (s *sRsa) init() {
	if s.privateKey != "" {
		// 解码将找到下一个PEM格式化块(证书，私钥，etc)
		block, _ := pem.Decode([]byte(s.privateKey))

		// 判断私钥类型
		//pkcs1
		if strings.Index(s.privateKey, "BEGIN RSA") > 0 {
			s.rsaPrivateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
		} else { //pkcs8
			privateKey, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
			s.rsaPrivateKey = privateKey.(*rsa.PrivateKey)
		}
	}

	if s.publicKey != "" {
		// 提取公钥
		block, _ := pem.Decode([]byte(s.publicKey))
		publicKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
		s.rsaPublicKey = publicKey.(*rsa.PublicKey)
	}
}

/**
 * 加密
 */
func (s *sRsa) Encrypt(data []byte) ([]byte, error) {
	// 客户端用服务端给的公钥进行数据加密 （客户端调用sdk）
	blockLength := s.rsaPublicKey.N.BitLen()/8 - 11
	if len(data) <= blockLength {
		return rsa.EncryptPKCS1v15(rand.Reader, s.rsaPublicKey, []byte(data))
	}

	buffer := bytes.NewBufferString("")

	pages := len(data) / blockLength

	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(data) {
				continue
			}
			end = len(data)
		}

		chunk, err := rsa.EncryptPKCS1v15(rand.Reader, s.rsaPublicKey, data[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

/**
 * 解密
 */
func (s *sRsa) Decrypt(secretData []byte) ([]byte, error) {
	blockLength := s.rsaPublicKey.N.BitLen() / 8
	// 服务端用给客户端公钥对应的私钥解密（如果公钥被拦截 无法识别用户 需要用到签名）
	if len(secretData) <= blockLength {
		return rsa.DecryptPKCS1v15(rand.Reader, s.rsaPrivateKey, secretData)
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

		chunk, err := rsa.DecryptPKCS1v15(rand.Reader, s.rsaPrivateKey, secretData[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

/**
 * 签名 客户端需要调用的接口，（客户本地的私钥和加密过的数据请求服务端）
 */
func (s *sRsa) Sign(data []byte, algorithmSign crypto.Hash) ([]byte, error) {
	hash := algorithmSign.New()
	hash.Write(data)
	//   用 用户私钥做签名
	sign, err := rsa.SignPKCS1v15(rand.Reader, s.rsaPrivateKey, algorithmSign, hash.Sum(nil))
	if err != nil {
		return nil, err
	}
	return sign, err
}

/**
 * 验签 服务端用客户端事先给服务端的公钥和签名 验签 （验签完成后才会去执行解密数据）
 */
func (s *sRsa) Verify(data []byte, sign []byte, algorithmSign crypto.Hash) bool {
	h := algorithmSign.New()
	h.Write(data)
	// 用c端用户公钥和签名做验签
	// 客户公钥需要保存在服务端
	return rsa.VerifyPKCS1v15(s.rsaPublicKey, algorithmSign, h.Sum(nil), sign) == nil
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

func main() {
	//content   := strings.Repeat("H", 244)+"e"
	//content   := strings.Repeat("H", 245)+"e"
	content := strings.Repeat("H", 24270) + "e"
	//privateKey, publicKey := NewRsa("", "").CreateKeys(1024)
	privateKey, publicKey := Rsa("", "").CreatePkcs8Keys(2048)
	fmt.Printf("公钥：%v\n私钥：%v\n", publicKey, privateKey)

	rsaObj := Rsa(publicKey, privateKey)
	// 加密
	secretData, err := rsaObj.Encrypt([]byte(content))
	if err != nil {
		fmt.Println(err)
	}
	// 解密
	plainData, err := rsaObj.Decrypt(secretData)
	if err != nil {
		fmt.Print(err)
	}

	data := []byte(strings.Repeat(content, 200))
	//sign,_ := rsaObj.Sign(data, crypto.SHA1)
	//verify := rsaObj.Verify(data, sign, crypto.SHA1)
	// 签名
	sign, _ := rsaObj.Sign(data, crypto.SHA256)
	// 验签
	verify := rsaObj.Verify(data, sign, crypto.SHA256)

	fmt.Printf(" 加密：%v\n 解密：%v\n 签名：%v\n 验签结果：%v\n",
		hex.EncodeToString(secretData),
		string(plainData),
		hex.EncodeToString(sign),
		verify,
	)
}
