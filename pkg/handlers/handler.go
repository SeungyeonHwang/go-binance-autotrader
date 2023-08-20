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
		{binance.SUB1_ACCOUNT, binance.SUB1_EMAIL, "Sub1"},
		// {binance.MASTER_ACCOUNT, binance.SUB2_EMAIL, "Sub2"},
		// {binance.MASTER_ACCOUNT, binance.SUB3_EMAIL, "Sub3"},
	}

	var resultBuilder strings.Builder

	for _, acc := range accounts {
		balance, err := binance.GetFuturesBalance(acc.accountType, config, acc.email)
		if err != nil {
			return "", err
		}

		resultBuilder.WriteString(fmt.Sprintf("Account     | %s\n", acc.label))
		resultBuilder.WriteString("--------------------\n")

		if acc.label == "Sub1" {
			unitPrice := int(float64(balance) * 0.05 * 15)
			resultBuilder.WriteString(fmt.Sprintf("Time        | 1H\n"))
			resultBuilder.WriteString(fmt.Sprintf("Leverage    | X15\n"))
			resultBuilder.WriteString(fmt.Sprintf("Unit Price  | %d\n", unitPrice))
			resultBuilder.WriteString(fmt.Sprintf("Balance     | %d\n", balance))
		} else {
			resultBuilder.WriteString(fmt.Sprintf("Balance     | %d\n", balance))
		}
		resultBuilder.WriteString("--------------------\n")
		resultBuilder.WriteString("\n")
	}

	return resultBuilder.String(), nil
}
