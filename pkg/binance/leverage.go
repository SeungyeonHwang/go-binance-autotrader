package binance

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

func changeLeverage(client *futures.Client, symbol string, leverage int) error {
	maxRetries := 12
	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		_, err := client.NewChangeLeverageService().Symbol(symbol).Leverage(leverage).Do(context.Background())
		if err == nil {
			return nil
		}

		leverage--
		if leverage <= 0 {
			break
		}
	}

	return fmt.Errorf("failed to change leverage after %d retries", maxRetries)
}

func changeMarginType(apiKey, secretKey, symbol, marginType string) error {
	endpoint := "/fapi/v1/marginType"
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("marginType", marginType)
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	sig := createHmac(params.Encode(), secretKey)
	params.Add("signature", sig)

	_, err := sendRequest(apiKey, "POST", baseURL+endpoint, params)
	return err
}
