package v1

import "github.com/gogf/gf/v2/frame/g"

type GethReq struct {
	g.Meta `path:"/geth" method:"get" summary:"测试调用eth" tags:"智能合约"`
}
type GethRes struct {
	Result string `json:"result" dc:"调用结果"`
}
