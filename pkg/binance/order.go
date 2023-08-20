package binance

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/SeungyeonHwang/go-binance-autotrader/config"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

func NewFuturesClient(config *config.Config, account string) (*futures.Client, error) {
	log.Printf("Initializing a new Futures client for account: %s", account)

	var apiKey, secretKey string
	switch strings.ToLower(account) {
	case MASTER_ACCOUNT:
		apiKey = config.MasterAPIKey
		secretKey = config.MasterSecretKey
	case SUB1_ACCOUNT:
		apiKey = config.Sub1APIKey
		secretKey = config.Sub1SecretKey
	default:
		return nil, fmt.Errorf("account %s not found in the configuration", account)
	}

	return binance.NewFuturesClient(apiKey, secretKey), nil
}

func PlaceFuturesMarketOrder(config *config.Config, account, symbol, positionSide string, leverage, amountInUSDT int) error {
	symbol = FormatSymbol(symbol)
	positionSide = ToUpper(positionSide)

	client, err := NewFuturesClient(config, account)
	if err != nil {
		return err
	}

	if err := changeLeverage(client.APIKey, client.SecretKey, symbol, leverage); err != nil {
		log.Printf("Failed to set leverage: %s", err)
		return err
	}
	if err := changeMarginType(client.APIKey, client.SecretKey, symbol, CROSSED); err != nil {
		log.Printf("Failed to set margin type: %s", err)
		return err
	}

	if err := setPositionSideMode(client.APIKey, client.SecretKey, true); err != nil {
		log.Printf("Failed to set position side mode to Hedge: %s", err)
	}

	price, err := getCurrentFuturesPrice(client, symbol)
	if err != nil {
		log.Printf("Failed to fetch the current price: %s", err)
		return err
	}

	quantity := float64(amountInUSDT) / price

	stepSize, err := getStepSizeForSymbol(client, symbol)
	if err != nil {
		log.Printf("Failed to fetch step size: %s", err)
		return err
	}
	trimmedQuantity := trimQuantity(quantity, stepSize)
	roundedQuantity := math.Round(trimmedQuantity*1e6) / 1e6

	var orderSide futures.SideType
	orderType := OPEN
	if orderType == OPEN {
		if positionSide == LONG {
			orderSide = futures.SideTypeBuy
		} else if positionSide == SHORT {
			orderSide = futures.SideTypeSell
		} else {
			return fmt.Errorf("invalid position side provided: %s", positionSide)
		}
	} else if orderType == CLOSE {
		if positionSide == SHORT {
			orderSide = futures.SideTypeBuy
			quantity = -quantity
		} else if positionSide == LONG {
			orderSide = futures.SideTypeSell
		} else {
			return fmt.Errorf("invalid position side provided: %s", positionSide)
		}
	} else {
		return fmt.Errorf("invalid order type provided: %s", orderType)
	}

	order, err := client.NewCreateOrderService().Symbol(symbol).
		Side(orderSide).PositionSide(futures.PositionSideType(positionSide)).Type(futures.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.6f", roundedQuantity)).Do(context.Background())
	if err != nil {
		log.Printf("Futures Market Order failed: %s", err)
		return err
	}
	log.Printf("Order ID: %d", order.OrderID)
	return nil
}
