package binance

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

type Position struct {
	Symbol           string `json:"symbol"`
	InitialMargin    string `json:"initialMargin"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	PositionAmt      string `json:"positionAmt"`
	EntryPrice       string `json:"entryPrice"`
	PositionSide     string `json:"positionSide"`
}

func getCurrentPosition(client *futures.Client, symbol string) (*futures.AccountPosition, error) {
	currentPosition, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error fetching current position: %v", err)
	}

	for _, pos := range currentPosition.Positions {
		if pos.Symbol == symbol {
			positionAmt, err := strconv.ParseFloat(pos.PositionAmt, 64)
			if err != nil {
				return nil, fmt.Errorf("error convert position amount: %v", err)
			}
			if positionAmt != 0.0 {
				return pos, nil
			}
		}
	}
	return nil, fmt.Errorf("no matching position found for symbol: %s", symbol)
}

func setPositionSideMode(apiKey, secretKey string, hedgeMode bool) error {
	endpoint := "/fapi/v1/positionSide/dual"
	params := url.Values{}
	params.Add("dualSidePosition", strconv.FormatBool(hedgeMode))
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	sig := createHmac(params.Encode(), secretKey)
	params.Add("signature", sig)

	_, err := sendRequest(apiKey, "POST", baseURL+endpoint, params)
	return err
}
