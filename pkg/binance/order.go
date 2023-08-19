package binance

import (
	"context"
	"fmt"
	"log"
	"math"

	binance "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

func NewFuturesClient(account string) (*futures.Client, error) {
	log.Printf("Initializing a new Futures client for account: %s", account)
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	if creds, ok := config.Binance[account]; ok {
		return binance.NewFuturesClient(creds.APIKey, creds.SecretKey), nil
	}
	return nil, fmt.Errorf("account %s not found in the configuration", account)
}

// (PlaceFuturesMarketOrder 함수 및 관련된 함수들)

func PlaceFuturesMarketOrder(account, symbol string, amountInUSDT float64, positionSide string) error {
	log.Printf("Placing Futures Market Order for account: %s, symbol: %s, amount: %f, position: %s", account, symbol, amountInUSDT, positionSide)
	client, err := NewFuturesClient(account)
	if err != nil {
		return err
	}

	// Set leverage to 1 and margin type to CROSSED
	if err := changeLeverage(client.APIKey, client.SecretKey, symbol, 1); err != nil {
		log.Printf("Failed to set leverage: %s", err)
		return err
	}
	if err := changeMarginType(client.APIKey, client.SecretKey, symbol, "CROSSED"); err != nil {
		log.Printf("Failed to set margin type: %s", err)
		return err
	}

	if err := setPositionSideMode(client.APIKey, client.SecretKey, true); err != nil {
		log.Printf("Failed to set position side mode to Hedge: %s", err)
	}

	// 현재 BTC의 가격을 가져옵니다.
	price, err := getCurrentFuturesPrice(client, symbol) // 여기서 client를 전달합니다.
	if err != nil {
		log.Printf("Failed to fetch the current price: %s", err)
		return err
	}

	// 주어진 USDT 양에 해당하는 BTC의 양을 계산합니다.
	quantity := amountInUSDT / price

	stepSize, err := getStepSizeForSymbol(account, symbol) // Added 'account' as a parameter
	if err != nil {
		log.Printf("Failed to fetch step size: %s", err)
		return err
	}
	trimmedQuantity := trimQuantity(quantity, stepSize)
	roundedQuantity := math.Round(trimmedQuantity*1e6) / 1e6 // Rounding to 6 decimal places

	var orderSide futures.SideType
	orderType := "OPEN" // orderType의 기본값을 "OPEN"으로 설정
	if orderType == "OPEN" {
		if positionSide == "LONG" {
			orderSide = futures.SideTypeBuy
		} else if positionSide == "SHORT" {
			orderSide = futures.SideTypeSell
		} else {
			return fmt.Errorf("invalid position side provided: %s", positionSide)
		}
	} else if orderType == "CLOSE" {
		if positionSide == "SHORT" {
			orderSide = futures.SideTypeBuy
			quantity = -quantity
		} else if positionSide == "LONG" {
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
