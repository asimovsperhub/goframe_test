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
func (s *sGethser) CommonEth(ctx context.Context, ethconnhost, calls string) (*tools.SAY, error) {
	// Dial connects a client to the given URL
	// ethconnhost 以太坊地址  calls 智能合约账户地址
	conn, err := ethclient.Dial(ethconnhost)
	if err != nil {
		panic("连接以太坊出错")
	}
	// 查询等操作可以连接返回后关闭连接
	// defer conn.Close()
	// 生成合约实例 , // NewERC20创建一个新的ERC20实例，绑定到一个特定的部署契约。
	gethObject, err := tools.NewSAY(common.HexToAddress(calls), conn)
	if err != nil {
		panic("创建合约对象出错")
	}
	return gethObject, nil
}

// todo 调用合约对象处理业务逻辑
func (s *sGethser) DoFunc(ctx context.Context) (string, error) {
	//
	gethObject, err := s.CommonEth(ctx, "http://localhost:8545", "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4")
	if err != nil {
		panic(fmt.Sprintf("创建合约对象出错 %s", err))
	}

	// 调用合约方法
	say, err := gethObject.Msg(nil)

	return say, err
}

//pargma  solidity^0.6.0;
//  contract Say{
//	string public Msg;
//	constructor() public{
//		Msg = "hello";
//}
//}
