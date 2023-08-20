package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SeungyeonHwang/go-binance-autotrader/config"
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/binance"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Echo   *echo.Echo
	Config *config.Config
}

func (h *Handler) CheckBalance(c echo.Context) error {
	val, err := fetchAllBalances(h.Config)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch balance")
	}
	return c.String(http.StatusOK, val)
}

func fetchAllBalances(config *config.Config) (string, error) {
	accounts := []struct {
		accountType string
		email       string
		label       string
	}{
		{binance.MASTER_ACCOUNT, "", "Master"},
		{binance.MASTER_ACCOUNT, binance.SUB1_EMAIL, "Sub1"},
		// {binance.MASTER_ACCOUNT, binance.SUB2_EMAIL, "Sub2"},
		// {binance.MASTER_ACCOUNT, binance.SUB3_EMAIL, "Sub3"},
	}
	balances := make([]int, len(accounts))
	resultStrings := make([]string, len(accounts))

	for i, acc := range accounts {
		balance, err := binance.GetFuturesBalance(acc.accountType, config, acc.email)
		if err != nil {
			return "", err
		}
		balances[i] = balance
		resultStrings[i] = fmt.Sprintf("%s Balance: %d", acc.label, balances[i])
	}

	result := strings.Join(resultStrings, "\n")
	return result, nil
}
