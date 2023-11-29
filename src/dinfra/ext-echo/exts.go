package extecho

import (
	"github.com/AlgerDu/go-dream/src/dinfra"
	"github.com/labstack/echo/v4"
)

func Bind[DstType any](c echo.Context) (*DstType, error) {
	data := new(DstType)
	err := c.Bind(data)

	return data, err
}

func SetResult(c echo.Context, code int) error {
	return c.JSON(200, &dinfra.Result{
		Code: code,
	})
}

func SetResultWith(c echo.Context, code int, respData any) error {
	return c.JSON(200, &dinfra.Result{
		Code: code,
		Data: respData,
	})
}
