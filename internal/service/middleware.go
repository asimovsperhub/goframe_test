package service

import "github.com/gogf/gf/v2/net/ghttp"

// 权限中间件
type middlewareService struct{}

var middleware = middlewareService{}

func Middleware() *middlewareService {
	return &middleware
}

func (s *middlewareService) Auth(r *ghttp.Request) {
	// GfJWTMiddleware gf jwt集成的中间件
	Auth().MiddlewareFunc()(r)
	r.Middleware.Next()
}
