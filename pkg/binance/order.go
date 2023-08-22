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

func PlaceFuturesMarketOrder(config *config.Config, account, symbol, positionSide string, leverage, amountInUSDT int, entry bool) error {
	symbol = FormatSymbol(symbol)
	positionSide = ToUpper(positionSide)

	client, err := NewFuturesClient(config, account)
	if err != nil {
		return err
	}

	// Entry == false
	if !entry {
		newAmount, err := roiValidation(client.APIKey, client.SecretKey, symbol, leverage, amountInUSDT)
		if err != nil {
			return err
		}
		amountInUSDT = newAmount
	} else {

		//Entry == true
		exists, err := positionExistsForSymbol(client.APIKey, client.SecretKey, symbol)
		if err != nil {
			return err
		}

		if entry && exists {
			newAmount, err := roiValidation(client.APIKey, client.SecretKey, symbol, leverage, amountInUSDT)
			if err != nil {
				return err
			}
			amountInUSDT = newAmount
		}
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
	_, err = client.NewCreateOrderService().Symbol(symbol).
		Side(orderSide).PositionSide(futures.PositionSideType(positionSide)).Type(futures.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.6f", roundedQuantity)).Do(context.Background())
	if err != nil {
		log.Printf("Futures Market Order failed: %s", err)
		return err
	}

	msg := ":bell: Order has been successfully created.\n"
	msg += "Account: " + account + "\n"
	msg += "Symbol: " + symbol + "\n"
	msg += "Position: " + positionSide + "\n"
	msg += "Leverage: " + fmt.Sprintf("x%d", leverage) + "\n"
	msg += "Order Amount: " + fmt.Sprintf("%d USDT", amountInUSDT)

	err = SendSlackNotification("https://hooks.slack.com/services/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", msg)
	if err != nil {
		log.Printf("Failed to send Slack notification: %s", err)
	}
	return nil
}

func roiValidation(apiKey, secretKey, targetSymbol string, leverage, amountInUSDT int) (int, error) {
	roi, err := getROIForSymbol(apiKey, secretKey, targetSymbol)
	if err != nil {
		log.Printf("Failed to get ROI: %s", err)
		return 0, err
	}

	// 1. roi = 0 -> error return
	if roi == 0 {
		return 0, fmt.Errorf("ROI is zero for symbol %s", targetSymbol)
	}

	// 2. roi -1 * leverage -> amountInUSDT (e.g -20% / leverage 15)
	if roi < -1.0*float64(leverage) {
		return amountInUSDT, nil
	}

	// 3. roi 1 * leverage -> amountInUSDT / 2 (e.g 30% / leverage 15)
	if roi > 1.0*float64(leverage) {
		return amountInUSDT / 2, nil
	}

	// 4. 그 외의 값 +-1 * leverage -> error return (e.g 1%/ leverage 15)
	if roi > -1.0*float64(leverage) && roi < 1.0*float64(leverage) {
		return 0, fmt.Errorf("ROI not within acceptable range for position with symbol %s", targetSymbol)
	}

	return 0, err
}
