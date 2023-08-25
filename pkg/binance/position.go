package binance

import (
	"net/url"
	"strconv"
	"time"
)

type Position struct {
	Symbol           string `json:"symbol"`
	InitialMargin    string `json:"initialMargin"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	PositionAmt      string `json:"positionAmt"`
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
