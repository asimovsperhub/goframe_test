package cmd

import (
	"context"
	"firstproject/internal/controller"
	"firstproject/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				//  结果返回后置中间价
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					controller.Hello, controller.Login, controller.Register,
				)
				// 需要鉴权的接口
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(service.Middleware().Auth)
					group.ALLMap(g.Map{
						"/user/info": controller.User.Info,
					})
				})

			})
			s.Run()
			return nil
		},
	}
)
