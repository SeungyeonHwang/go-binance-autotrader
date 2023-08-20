package router

import (
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/handlers"
	"github.com/labstack/echo/v4"
)

func SetUp(e *echo.Echo, h *handlers.Handler) {
	swingAPI := e.Group("swing")
	{
		swingAPI.GET("/balance", h.CheckBalance)
	}
}