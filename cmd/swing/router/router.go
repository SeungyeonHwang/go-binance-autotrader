package router

import (
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/handlers"
	"github.com/labstack/echo/v4"
)

func SetUp(e *echo.Echo, h *handlers.Handler) {
	swingAPI := e.Group("swing")
	{
		swingAPI.GET("/balance", h.CheckBalance)
		swingAPI.GET("/history", h.CheckHistory)
		swingAPI.GET("/position", h.CheckPosition)
		swingAPI.POST("/webhook-order", h.WebhookOrder)
		swingAPI.POST("/sltp-all", h.SetStopLossTakeProfitALL)
		swingAPI.POST("/sltp-partial", h.SetStopLossTakeProfitPartial)
		swingAPI.POST("/db-clear", h.DBClear)
	}
}
