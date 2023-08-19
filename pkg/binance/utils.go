package binance

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
)

func trimQuantity(quantity, stepSize float64) float64 {
	trimmedQuantity := math.Round(quantity/stepSize) * stepSize
	return trimmedQuantity
}

func getStepSizeForSymbol(account, symbol string) (float64, error) {
	client, err := NewFuturesClient(account)
	if err != nil {
		return 0, err
	}
	// Fetch the exchange info
	info, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return 0, err
	}
	for _, s := range info.Symbols {
		if s.Symbol == symbol {
			// Parse the stepSize from the symbol's filter.
			// Assuming the filter type for step size is "LOT_SIZE".
			for _, f := range s.Filters {
				if f["filterType"] == "LOT_SIZE" {
					return strconv.ParseFloat(f["stepSize"].(string), 64)
				}
			}
		}
	}
	return 0, fmt.Errorf("step size for symbol %s not found", symbol)
}

func getCurrentFuturesPrice(client *futures.Client, symbol string) (float64, error) {
	log.Printf("Fetching current price for symbol: %s", symbol)
	prices, err := client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return 0, err
	}
	for _, p := range prices {
		if p.Symbol == symbol {
			return strconv.ParseFloat(p.Price, 64)
		}
	}
	return 0, fmt.Errorf("price for symbol %s not found", symbol)
}
