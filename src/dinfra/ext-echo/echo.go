package extecho

import (
	"github.com/AlgerDu/go-dream/src/dinfra"
	"github.com/labstack/echo/v4"
)

func New(
	logger dinfra.Logger,
) (*dinfra.HttpServer, error) {

	e := echo.New()
	e.HideBanner = true
	e.Debug = true

	logger.Info("create echo http server")

	return &dinfra.HttpServer{
		Echo: e,
	}, nil
}
