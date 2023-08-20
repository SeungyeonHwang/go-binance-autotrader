package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/binance"
	"github.com/labstack/echo/v4"
)

type Config struct {
	MasterAPIKey    string
	MasterSecretKey string
	Sub1APIKey      string
	Sub1SecretKey   string
}

type Handler struct {
	Echo   *echo.Echo
	Config *Config
}

func (h *Handler) CheckBalance(c echo.Context) error {
	// balance, err := fetchAllBalances()
	// if err != nil {
	// 	return c.String(http.StatusInternalServerError, "Failed to fetch balance")
	// }
	log.Printf("MasterAPIKey: %s", h.Config.MasterAPIKey)
	log.Printf("MasterSecretKey: %s", h.Config.MasterSecretKey)
	log.Printf("Sub1APIKey: %s", h.Config.Sub1APIKey)
	log.Printf("Sub1SecretKey: %s", h.Config.Sub1SecretKey)
	return c.String(http.StatusOK, "helloworld")
}

func fetchAllBalances() (string, error) {
	masterBalance, err := binance.GetFuturesBalance(binance.MASTER_ACCOUNT)
	if err != nil {
		return "", err
	}

	// sub1Balance, err := binance.GetFuturesBalance(binance.SUB1_ACCOUNT, binance.SUB1_EMAIL)
	// if err != nil {
	// 	return "", err
	// }

	// sub2Balance, err := binance.GetFuturesBalance(binance.SUB2_ACCOUNT, binance.SUB2_EMAIL)
	// if err != nil {
	// 	return "", err
	// }

	// sub3Balance, err := binance.GetFuturesBalance(binance.SUB3_ACCOUNT, binance.SUB3_EMAIL)
	// if err != nil {
	// 	return "", err
	// }

	result := fmt.Sprintf("Master Balance: %d", masterBalance)
	// result := fmt.Sprintf("Master Balance: %d\nSub1: %d\nSub2: %d\nSub3: %d", masterBalance, sub1Balance, sub2Balance, sub3Balance)

	return result, nil
}
