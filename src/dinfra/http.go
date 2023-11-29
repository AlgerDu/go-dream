package dinfra

import "github.com/labstack/echo/v4"

// 针对 echo 的一些简单封装

const (
	HttpHeader_XRequestID = "X-Request-Id"
)

type (
	// 自定义中间件接口，方便实现自定义的中间件（ echo 原始的中间件实现方式，看起来比较难看，仅此而已）
	HttpMiddleware interface {
		Handle(context echo.Context, next echo.HandlerFunc) error
	}

	// 简单包装的一个 http 服务
	HttpServer struct {
		*echo.Echo
	}
)

// 将自定义 HttpMiddleware 接口转换为 echo 可用的 MiddlewareFunc
func Http_MiddlewareToEchoFunc(middleware HttpMiddleware) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			return middleware.Handle(context, next)
		}
	}
}

func (server *HttpServer) Start(address string) error {
	return server.Echo.Start(address)
}
