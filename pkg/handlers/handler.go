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
		return c.String(http.StatusInternalServerError, "Error fetching history.")
	}

	if c.Path() == "/swing/history" {
		return c.String(http.StatusOK, val)
	}

	return c.NoContent(http.StatusOK)
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
func (h *Handler) WebhookOrder(c echo.Context) error {
	orderReq := NewTradingViewPayload()
	if err := c.Bind(&orderReq); err != nil {
		return c.String(http.StatusBadRequest, "Failed to parse request body")
	}

	err := binance.PlaceFuturesMarketOrder(h.Config, orderReq.Account, orderReq.Symbol, orderReq.PositionSide, orderReq.Amount, orderReq.Entry)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to place order: %s", err.Error()))
	}
	return c.String(http.StatusOK, "Order placed successfully!")
}

func (h *Handler) SetStopLossTakeProfitAll(c echo.Context) error {
	orderReq := new(AllStopLossTakeProfitPayload)
	if err := c.Bind(orderReq); err != nil {
		return c.String(http.StatusBadRequest, "Failed to parse request body")
	}

	err := binance.PlaceALLStopLossTakeProfitOrder(h.Config, orderReq.Account, orderReq.Symbol, orderReq.TP, orderReq.SL)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to set stop loss and take profit: %s", err.Error()))
	}
	return c.String(http.StatusOK, "SL/TP(ALL) set successfully!")
}

func (h *Handler) SetStopLossTakeProfitPartial(c echo.Context) error {
	orderReq := new(PartialTakeProfitPayload)
	if err := c.Bind(orderReq); err != nil {
		return c.String(http.StatusBadRequest, "Failed to parse request body")
	}

	details := make(map[string]interface{})
	details["account"] = orderReq.Account
	details["symbol"] = orderReq.Symbol
	details["positionSide"] = orderReq.PositionSide
	if orderReq.TP != nil {
		tpMap := make(map[string]interface{})
		tpMap["price"] = orderReq.TP.Price
		tpMap["quantity"] = orderReq.TP.Quantity
		details["tp"] = tpMap
	}

	err := binance.PlacePartialTakeProfitOrder(h.Config, details)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to set partial take profit: %s", err.Error()))
	}
	return c.String(http.StatusOK, "TP(Partial) set successfully!")
}

// Close order
func (h *Handler) CloseOrder(c echo.Context) error {
	orderReq := new(CloseOrderPayload)
	if err := c.Bind(orderReq); err != nil {
		return c.String(http.StatusBadRequest, "Failed to parse request body")
	}

	err := binance.PlaceFuturesMarketCloseOrder(h.Config, orderReq.Account, orderReq.Symbol, orderReq.Close)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to place close order: %s", err.Error()))
	}
	return c.String(http.StatusOK, "Close Order placed successfully!")
}

func (h *Handler) CloseAllOrder(c echo.Context) error {
	orderReq := new(CloseAllOrderPayload)
	if err := c.Bind(orderReq); err != nil {
		return c.String(http.StatusBadRequest, "Failed to parse request body")
	}

	err := binance.PlaceFuturesMarketCloseAllOrder(h.Config, orderReq.Account)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to place close all order: %s", err.Error()))
	}
	return c.String(http.StatusOK, "Close All Order placed successfully!")
}

// db-clear
func (h *Handler) DBClear(c echo.Context) error {
	err := binance.DBClear(binance.BUCKET_NAME, binance.DB_NAME)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to clear db")
	}
	return c.String(http.StatusOK, "DB cleared")
}
