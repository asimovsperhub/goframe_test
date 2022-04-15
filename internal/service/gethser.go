package service

import (
	"context"
	"firstproject/tools"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type sGethser struct{}

// 管理职能合约服务
func Gethser() *sGethser {
	return &sGethser{}
}

// 创建合约对象
func (s *sGethser) CommonEth(ctx context.Context, ethconnhost, calls string) interface{} {
	// Dial connects a client to the given URL
	// ethconnhost 以太坊地址  calls 智能合约账户地址
	conn, err := ethclient.Dial(ethconnhost)
	if err != nil {
		panic("连接以太坊出错")
	}
	// 查询等操作可以连接返回后关闭连接
	// defer conn.Close()
	// 生成合约实例 , NewCallOpts creates a new option set for contract calls.
	gethObject, err := tools.NewERC20Caller(common.HexToAddress(calls), conn)
	if err != nil {
		panic("创建合约对象出错")
	}
	return gethObject
}

// todo 调用合约对象处理业务逻辑
func (s *sGethser) DoFunc(ctx context.Context) error {
	//
	s.CommonEth(ctx, "", "")
	return nil
}
