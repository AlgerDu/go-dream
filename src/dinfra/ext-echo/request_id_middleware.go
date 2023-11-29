package extecho

import (
	"github.com/AlgerDu/go-dream/src/dinfra"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	RequestIDMiddleware struct {
		logger dinfra.Logger
	}
)

func NewRequestIDMiddleware(
	logger dinfra.Logger,
) *RequestIDMiddleware {

	return &RequestIDMiddleware{
		logger: logger.WithField(dinfra.LogField_Source, "RequestIDMiddleware"),
	}
}

func (middleware *RequestIDMiddleware) Handle(c echo.Context, next echo.HandlerFunc) error {

	req := c.Request()
	res := c.Response()
	rid := req.Header.Get(dinfra.HttpHeader_XRequestID)
	if rid == "" {
		rid = uuid.NewString()
		middleware.logger.WithField(dinfra.HttpHeader_XRequestID, rid).Info("generate htttp request id")
	}

	res.Header().Set(dinfra.HttpHeader_XRequestID, rid)
	SetRequestID(c, rid)

	return next(c)
}
