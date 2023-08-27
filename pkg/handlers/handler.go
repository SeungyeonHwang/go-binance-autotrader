package handlers

import (
	"fmt"
	"net/http"

	"github.com/SeungyeonHwang/go-binance-autotrader/config"
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/binance"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Echo   *echo.Echo
	Config *config.Config
}

// balance
func (h *Handler) CheckBalance(c echo.Context) error {
	val, err := binance.FetchAllBalances(h.Config)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch balance")
	}
	return c.String(http.StatusOK, val)
}

// balance
func (h *Handler) CheckHistory(c echo.Context) error {
	val, err := binance.FetchAllHistory(h.Config, binance.BUCKET_NAME, binance.DB_NAME)
	if err != nil {
		return c.String(http.StatusInternalServerError, val)
	}
	return c.String(http.StatusOK, val)
}

// position
func (h *Handler) CheckPosition(c echo.Context) error {
	val, err := binance.FetchAllPositions(h.Config)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch position")
	}
	return c.String(http.StatusOK, val)
}

// webhook trigger order
// https://x8oktqy9c1.execute-api.ap-northeast-1.amazonaws.com/Prod/swing/webhook-order
//
//		{
//		  "account": "master" or "sub1",
//		  "symbol": "{{ticker}}",
//		  "positionSide": "long",
//		  "leverage":15,
//		  "amount": 30(unit price),
//	   "entry":true(default false)
//		}
func (h *Handler) WebhookOrder(c echo.Context) error {
	orderReq := new(TradingViewPayload)
	if err := c.Bind(orderReq); err != nil {
		return c.String(http.StatusBadRequest, "Failed to parse request body")
	}

	err := binance.PlaceFuturesMarketOrder(h.Config, orderReq.Account, orderReq.Symbol, orderReq.PositionSide, orderReq.Leverage, orderReq.Amount, orderReq.Entry)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to place order: %s", err.Error()))
	}
	return c.String(http.StatusOK, "Order placed successfully!")
}

func (h *Handler) SetStopLossTakeProfitALL(c echo.Context) error {
	orderReq := new(StopLossTakeProfitPayload)
	if err := c.Bind(orderReq); err != nil {
		return c.String(http.StatusBadRequest, "Failed to parse request body")
	}

	err := binance.PlaceStopLossTakeProfitALLOrder(h.Config, orderReq.Account, orderReq.Symbol, orderReq.PositionSide, orderReq.TP, orderReq.SL)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to set stop loss and take profit: %s", err.Error()))
	}
	return c.String(http.StatusOK, "SL/TP(ALL) set successfully!")
}

// db-clear
func (h *Handler) DBClear(c echo.Context) error {
	err := binance.DBClear(binance.BUCKET_NAME, binance.DB_NAME)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to clear db")
	}
	return c.String(http.StatusOK, "DB cleared")
}
