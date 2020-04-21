package middleware

import (
	"github.com/dshechao/gen-api/gen"
	"github.com/dshechao/gen-api/gen/models"
	"github.com/dshechao/gen-api/middleware"
	"github.com/labstack/echo"
)

func Yaag() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			if !gen.IsOn() {
				return next(c)
			}

			apiCall := models.ApiCall{}
			writer := middleware.NewResponseRecorder(c.Response().Writer)
			c.Response().Writer = writer
			middleware.Before(&apiCall, c.Request())
			err := next(c)
			middleware.After(&apiCall, writer, c.Request())
			return err
		})
	})
}
