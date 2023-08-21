package middlewares

import (
	"fmt"
	"lift-fitness-gym/app/pkg/applog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)
var logger = applog.Get()
var LoggerMiddleware = middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogURI:      true,
	LogStatus:   true,
	LogMethod:   true,
	LogLatency:  true,
	LogRemoteIP: true,
	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
		contentType := c.Request().Header.Get("content-type")
		if v.Status >= 400 {
			logger.Error("Error", zap.String("URI", v.URI), zap.Int("status", v.Status), zap.String("duration", fmt.Sprint(v.Latency.Milliseconds(), " ms")), zap.String("method", v.Method), zap.String("contentType", contentType), zap.String("IP", v.RemoteIP))
		}
		if v.Status >= 300 && v.Status < 400 {
			logger.Info("Redirect", zap.String("URI", v.URI), zap.Int("status", v.Status), zap.String("duration", fmt.Sprint(v.Latency.Milliseconds(), " ms")), zap.String("method", v.Method), zap.String("contentType", contentType), zap.String("IP", v.RemoteIP))
		}
		if v.Status >= 200 && v.Status < 300 {
			logger.Info("Request", zap.String("URI", v.URI), zap.Int("status", v.Status), zap.String("duration", fmt.Sprint(v.Latency.Milliseconds(), " ms")), zap.String("method", v.Method), zap.String("contentType", contentType), zap.String("IP", v.RemoteIP))
		}

		return nil
	},
})