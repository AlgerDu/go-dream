package extecho

import (
	"github.com/AlgerDu/go-dream/src/dinfra"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RequestIDMiddleware struct {
	logger dinfra.Logger
}

func NewRequestIDMiddleware(
	logger dinfra.Logger,
) *RequestIDMiddleware {

	return &RequestIDMiddleware{
		logger: logger,
	}
}

func (middleware *RequestIDMiddleware) Handle(context echo.Context, next echo.HandlerFunc) error {

	requestID := context.Request().Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.NewString()
	}

	SetRequestID(context, requestID)

	return next(context)
}
