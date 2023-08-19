package binance

import (
	"net/url"
	"strconv"
	"time"
)

func changeLeverage(apiKey, secretKey, symbol string, leverage int) error {
	endpoint := "/fapi/v1/leverage"
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("leverage", strconv.Itoa(leverage))
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	sig := createHmac(params.Encode(), secretKey)
	params.Add("signature", sig)

	_, err := sendRequest(apiKey, "POST", baseURL+endpoint, params)
	return err
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
