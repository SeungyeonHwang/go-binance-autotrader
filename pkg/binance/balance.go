package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetFuturesBalance(account string, subAccountEmail ...string) (int, error) {
	// config, err := getConfig()
	// if err != nil {
	// 	return 0, fmt.Errorf("Could not fetch configuration: %v", err)
	// }
	// config, err := getConfig()
	// 	if err != nil {
	// 		return 0, fmt.Errorf("Could not fetch configuration: %v", err)
	// 	}

	// 	if account == "master" || len(subAccountEmail) == 0 {
	// 		return fetchBalance(config.Binance[account])
	// 	} else {
	// 		return fetchSubAccountBalance(config.Binance["master"], subAccountEmail[0])
	// 	}
	// }

	// func fetchBalance(credentials struct {
	// 	APIKey    string `yaml:"api_key"`
	// 	SecretKey string `yaml:"secret_key"`
	// }) (int, error) {
	// 	url := baseURL + "/fapi/v2/balance"
	// 	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	// 	queryString := "timestamp=" + timestamp
	// 	signature := createHmac(queryString, credentials.SecretKey)

	// 	req, err := http.NewRequest("GET", url+"?"+queryString+"&signature="+signature, nil)
	// 	if err != nil {
	// 		log.Fatalf("Could not create request: %v", err)
	// 	}
	// 	req.Header.Add("X-MBX-APIKEY", credentials.APIKey)

	// 	resp, err := http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		log.Fatalf("Request failed: %v", err)
	// 	}
	// 	defer resp.Body.Close()

	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Fatalf("Failed to read response body: %v", err)
	// 	}

	// 	var balanceData []map[string]interface{}
	// 	if err := json.Unmarshal(body, &balanceData); err != nil {
	// 		log.Fatalf("JSON unmarshalling failed: %s", err)
	// 	}

	// 	for _, asset := range balanceData {
	// 		if asset["asset"] == "USDT" {
	// 			balanceFloat, err := strconv.ParseFloat(asset["balance"].(string), 64)
	// 			if err != nil {
	// 				return 0, fmt.Errorf("Failed to parse balance to float: %v", err)
	// 			}
	// 			balanceInt := int(balanceFloat)
	// 			return balanceInt, nil
	// 		}
	// 	}
	return 0, fmt.Errorf("USDT not found in the response")
}

func fetchSubAccountBalance(masterCredentials struct {
	APIKey    string `yaml:"api_key"`
	SecretKey string `yaml:"secret_key"`
}, subAccountEmail string) (int, error) {
	// 엔드포인트
	url := "https://api.binance.com/sapi/v1/sub-account/futures/account"
	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	queryString := "timestamp=" + timestamp + "&email=" + subAccountEmail
	signature := createHmac(queryString, masterCredentials.SecretKey)

	req, err := http.NewRequest("GET", url+"?"+queryString+"&signature="+signature, nil)
	if err != nil {
		log.Fatalf("Could not create request: %v", err)
	}
	req.Header.Add("X-MBX-APIKEY", masterCredentials.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// 응답 형식에 맞게 구조체 정의
	type AssetInfo struct {
		Asset                  string  `json:"asset"`
		WalletBalance          float64 `json:"walletBalance,string"`
		UnrealizedProfit       float64 `json:"unrealizedProfit,string"`
		MarginBalance          float64 `json:"marginBalance,string"`
		MaintenanceMargin      float64 `json:"maintenanceMargin,string"`
		InitialMargin          float64 `json:"initialMargin,string"`
		PositionInitialMargin  float64 `json:"positionInitialMargin,string"`
		OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"`
		MaxWithdrawAmount      float64 `json:"maxWithdrawAmount,string"`
	}
	type ResponseData struct {
		Assets []AssetInfo `json:"assets"`
	}
	var responseData ResponseData

	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Fatalf("JSON unmarshalling failed: %s", err)
	}

	for _, asset := range responseData.Assets {
		if asset.Asset == "USDT" {
			return int(asset.WalletBalance), nil
		}
	}
	return 0, fmt.Errorf("USDT not found in the response")
}
