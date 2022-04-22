package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

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
 * 签名 客户端需要调用的接口，（客户本地的私钥和加密过的数据请求服务端）
 */
func (s *sRsa) Sign(data []byte, algorithmSign crypto.Hash) ([]byte, error) {
	hash := algorithmSign.New()
	hash.Write(data)
	//   用户私钥做签名
	sign, err := rsa.SignPKCS1v15(rand.Reader, s.rsaPrivateKey, algorithmSign, hash.Sum(nil))
	if err != nil {
		return nil, err
	}
	return sign, err
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

func POST(test_url string, data url.Values) {
	resp, err := http.PostForm(test_url, data)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	//fmt.Println(resp.Header)
}

func main() {
	// 测试解密
	// 需要带服务端给的 publicKey 请求get_key接口获取
	publicKey := "-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArDc0i4ZWuEqpbPvrHi6e\\nxLJIDHnreYMuJfeRPJq0Az0E4tsrlwQXy0hsmk5nIGwrdEBKJ30Er75O4V3SfwcU\\nFZYqnZ4vdqCIkTzwTYWGPOXb1mmbs8KuJzmmF7YfhwLH+DAbQOvJYteZqTWQiQtq\\njcw50rm8x3ifzB4QhdetJHQfdA8OaBXdzHeAAD8gtPilOi66n1lrR1CkS6uVpquP\\nRq0eVlJnJ56IfUVKsrP9bVH0dFfWcL7qXO8zJc3oJ2V45SobvXksGJbq6h3kGc99\\nU+k/bm9uppaRg/zJJ0IM26ioRQ0yl5NmIwfmVfSXlvuku47Fi5i7RMOWjW/ZJBM5\\n1wIDAQAB\\n-----END PUBLIC KEY-----\\n"
	// 私钥会在数据库查询当前用户的
	pub := strings.ReplaceAll(publicKey, "\\n", "\n")
	rsaObj := Rsa(pub, "")
	// 加密内容
	secretData, err := rsaObj.Encrypt([]byte("test"))
	if err != nil {
		println(err)
	}
	data1 := make(url.Values)
	data1["name"] = []string{"asimov123"}
	data1["secret"] = []string{string(secretData)}
	data1["publickey"] = []string{pub}
	test_url := "http://127.0.0.1:8000/decrypt"
	POST(test_url, data1) //{"code":0,"message":"","data":{"result":"Decrypt success:test"}}

	// 测试签名
	// 先push本地公钥
	public := "-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5ZPKd/oNXOp9Mhr6DnWd\\na1DEw+YjVVKKRY8edibRlvGhkAx3rKeLbNO1ZYhyYza8YPQEDjrAnZjD/prfh+HW\\nhtA39LLJjDXAjTlGb/2OsmxRAtmPALoPBbFHladV2O6ZBWI+p1sY+ODXRM1NRpKB\\n7aJut6Kg+wo/B0DwNioWPlqqH2OHxniwtTT9ZUCYQMXjCUmoa5a9LZBoFzvXmgnI\\nkeFpTeGJb0fJoZ7YUoZ+mUE5NhzXiLspcO9ZzaKkJlUSkYuJ0D5PR9s3tzIniifb\\nYrfPBWFyLIYpLEW4tvDmR3N993GRXyERJemVsZjBbnLr5mLVIYas9o0rmhEU5oMI\\nRQIDAQAB\\n-----END PUBLIC KEY-----\\n"
	data2 := make(url.Values)
	data2["name"] = []string{"asimov123"}
	data2["public"] = []string{strings.ReplaceAll(public, "\\n", "\n")}
	test_url1 := "http://127.0.0.1:8000/push_key"
	POST(test_url1, data2)

	// 签署信息
	data := []byte(strings.Repeat("test", 200))
	// 需要用户本地私钥
	rsaPrivateKey := "-----BEGIN PRIVATE KEY-----\\nMIIEvAIBADALBgkqhkiG9w0BAQEEggSoMIIEpAIBAAKCAQEA5ZPKd/oNXOp9Mhr6\\nDnWda1DEw+YjVVKKRY8edibRlvGhkAx3rKeLbNO1ZYhyYza8YPQEDjrAnZjD/prf\\nh+HWhtA39LLJjDXAjTlGb/2OsmxRAtmPALoPBbFHladV2O6ZBWI+p1sY+ODXRM1N\\nRpKB7aJut6Kg+wo/B0DwNioWPlqqH2OHxniwtTT9ZUCYQMXjCUmoa5a9LZBoFzvX\\nmgnIkeFpTeGJb0fJoZ7YUoZ+mUE5NhzXiLspcO9ZzaKkJlUSkYuJ0D5PR9s3tzIn\\niifbYrfPBWFyLIYpLEW4tvDmR3N993GRXyERJemVsZjBbnLr5mLVIYas9o0rmhEU\\n5oMIRQIDAQABAoIBAGNpgwQ/EGhK1hnLWrrGLXuaBwp5bpV035FNbzhkiN+fFIIH\\nFA98ocBnUKZ91mKmAh7Nq6/puxzDWSO4NtFldvr70S8x+FqxsAa3ZYv7NT6H7vCX\\n+veqmfSyFrh0NJVyhGqzZ0QbC45B9pXBfRPxPzgC3YTBdIognrhqY1phES7AS9hQ\\ntOfDD5B5jiAaIWXuDX9aXBJWT+NV0p59l+/xRTVen9bSVU3rl3rM1T7w5/Ax2l4G\\nvv603m4MIEMqXvufhvErmwBMNrpLIEjvy6zmfTqXoOCWoTIl0GNLQNaMulGvAf6T\\nL6xJ+kJ5BMxaeZVm4GU6QBaFpJ2gMVRekfXJHcUCgYEA6tRXywQcirvKgrb5NoGz\\nYf9Q5haXzdDbJHs6HbSySfS69fGurhaOJzywPUxa907OGbJ3zFMEpA3Z/R5Lzt5j\\nKBYXMOC7J9bN8jsPiR6RAtFIKzYDSX2/Xln3jydNLmVfvGgU4GqZmsacigNcG2hi\\nZ6qUj2j0qthPoQ7nu6YUszsCgYEA+kY7mXc+k0q++KGzm7syxHcmtP44w9llQmp0\\nLXlVm55scyYN9fbsTOIMDtNYgc4nlAcDT3vcDh/YxFAOZ3ZhFSc0aCCejxOUgV+F\\nSpWIxb22FCcN0m/yZx03hzB02P6dqxWEYVxcdushESchbvicRJUNIElyMylbDrg6\\n6Lk8en8CgYBekRSp1QYJeIadDUJfCOxMUp0pi3+miq01i8pjnBkQX1XLJYDK6ppk\\ngrQWe2FGpp2pC43i4qvDxTA8Fq9Ap54Wzo6YSGgWKxLUsaQX/A85qz386Mt6FQGz\\n5VckdxdFz9016lQ966/f/IudqKy2/NpkFPWuqv2cr2+h1HbNwpwjcQKBgQD4QBg4\\nNub0FW1ulH7jF4HZDVNwrsbBxe9CPPP2c2duUGvEoFeyxfZIoORTBGLDhykNFRO8\\nkOCLhh1vRPW0vOC5qcS7ELgWtdZVqdk+TSt48aAdR0vXlEF+9KUyzObqo0zj+hjw\\ntjvlnX+UUxs/xwzCnpKBlzjW9Mukwytz0uHhowKBgQCNk8zbbWgU8MXiylLQxcdu\\n45AFt7b8UW9H3y9I7ds/l8mQF9aHvAykpp8OSdNP5dASlYaFzq+2RaFMTUcZClxa\\ntXSxuPLAnUXDok5oP2A6SCX1LA9vJ2kxpLprh4BHaKwevbmx6qtyRHTJyc+V9L5a\\nTiD4zVzrW9Ez13heA3P4RQ==\\n-----END PRIVATE KEY-----\\n"
	rsaObj1 := Rsa("", strings.ReplaceAll(rsaPrivateKey, "\\n", "\n"))
	sign, _ := rsaObj1.Sign(data, crypto.SHA256)
	data3 := make(url.Values)
	data3["name"] = []string{"asimov123"}
	data3["sign"] = []string{string(sign)}
	data3["data"] = []string{string(data)}
	test_url2 := "http://127.0.0.1:8000/verify"
	POST(test_url2, data3) //{"code":0,"message":"","data":{"result":true}}

}
