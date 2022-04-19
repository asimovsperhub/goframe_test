package service

import (
	"context"
	"firstproject/tools"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type sGethser struct{}

// 管理智能合约服务
func Gethser() *sGethser {
	return &sGethser{}
}

// 创建合约对象
func (s *sGethser) CommonEth(ctx context.Context, ethconnhost, calls string) (*tools.CONTF, error) {
	// Dial connects a client to the given URL
	// ethconnhost 以太坊地址  calls 智能合约账户地址
	conn, err := ethclient.Dial(ethconnhost)
	if err != nil {
		panic("连接以太坊出错")
	}
	// 查询等操作可以连接返回后关闭连接
	// defer conn.Close()
	// 生成合约实例 , // NewCONTF创建一个新实例，绑定到一个特定的部署契约。
	gethObject, err := tools.NewCONTF(common.HexToAddress(calls), conn)
	if err != nil {
		panic("创建合约对象出错")
	}
	return gethObject, nil
}

// todo 调用合约对象处理业务逻辑
func (s *sGethser) DoFunc(ctx context.Context) (string, error) {
	//
	gethObject, err := s.CommonEth(ctx, "https://ropsten.infura.io/v3/0xd07bdb622a7e9d519a17c4c097bc479012761880", "0x2B52E4cA91F7d7a10152f42DF24139477e989b6B")
	if err != nil {
		panic(fmt.Sprintf("创建合约对象出错 %s", err))
	}

	//// 调用合约方法
	//gethObject.GetRandomWord(nil, nil)
	//_, err = gethObject.SetWord(nil, "s", "c")
	//if err != nil {
	//	panic(fmt.Sprintf("Call 合约出错 %s", err))
	//}
	return fmt.Sprintf("%s", gethObject), err
}

//pargma  solidity^0.6.0;
//  contract Say{
//	string public Msg;
//	constructor() public{
//		Msg = "hello";
//}
//}
