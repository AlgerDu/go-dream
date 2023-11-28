package extecho

import "github.com/labstack/echo/v4"

func SetRequestID(context echo.Context, value string) {
	context.Set("requestID", value)
}

func GetRequestID(context echo.Context) string {
	return context.Get("requestID").(string)
}
