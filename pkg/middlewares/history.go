package middlewares

import (
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/handlers"
	"github.com/labstack/echo/v4"
)

func HistoryMiddleware(h *handlers.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() != "/swing/history" && c.Path() != "/swing/db-clear" {
				if err := h.CheckHistory(c); err != nil {
					return err
				}
			}

			return next(c)
		}
	}
}
