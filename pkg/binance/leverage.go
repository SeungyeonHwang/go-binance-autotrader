package binance

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func changeLeverage(apiKey, secretKey, symbol string, leverage int) error {
	endpoint := "/fapi/v1/leverage"
	retryCount := 0
	maxRetries := 5

	for retryCount < maxRetries {
		params := url.Values{}
		params.Add("symbol", symbol)
		params.Add("leverage", strconv.Itoa(leverage))
		params.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
		sig := createHmac(params.Encode(), secretKey)
		params.Add("signature", sig)

		_, err := sendRequest(apiKey, "POST", baseURL+endpoint, params)
		if err == nil {
			return nil
		}

		leverage--
		retryCount++
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
