package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SeungyeonHwang/go-binance-autotrader/config"
)

func GetFuturesBalance(account string, config *config.Config, subAccountEmail ...string) (int, error) {
	if account == "master" || len(subAccountEmail) == 0 {
		return fetchBalance(config.MasterAPIKey, config.MasterSecretKey)
	} else {
		// You can further modify this line if there's an actual function to fetch sub account balance.
		// return fetchSubAccountBalance(config.MasterAPIKey, config.MasterSecretKey, subAccountEmail[0])
		return 0, nil
	}
}

func fetchBalance(apiKey string, secretKey string) (int, error) {
	url := baseURL + "/fapi/v2/balance"
	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	queryString := "timestamp=" + timestamp
	signature := createHmac(queryString, secretKey)

	req, err := http.NewRequest("GET", url+"?"+queryString+"&signature="+signature, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return 0, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Add("X-MBX-APIKEY", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return 0, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	var balanceData []map[string]interface{}
	if err := json.Unmarshal(body, &balanceData); err != nil {
		log.Printf("JSON unmarshalling error. Body: %s, Error: %v", string(body), err)
		return 0, fmt.Errorf("JSON unmarshalling failed: %s", err)
	}

	for _, assetData := range balanceData {
		if asset, exists := assetData["asset"]; exists && asset == "USDT" {
			balanceFloat, err := strconv.ParseFloat(assetData["balance"].(string), 64)
			if err != nil {
				log.Printf("Error parsing balance to float: %v", err)
				return 0, fmt.Errorf("failed to parse balance to float: %v", err)
			}
			return int(balanceFloat), nil
		}
	}

	log.Println("USDT not found in the response")
	return 0, fmt.Errorf("USDT not found in the response")
}

//TODO : balance for sub account
// func fetchSubAccountBalance(masterCredentials struct {
// 	APIKey    string `yaml:"api_key"`
// 	SecretKey string `yaml:"secret_key"`
// }, subAccountEmail string) (int, error) {
// 	// 엔드포인트
// 	url := "https://api.binance.com/sapi/v1/sub-account/futures/account"
// 	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
// 	queryString := "timestamp=" + timestamp + "&email=" + subAccountEmail
// 	signature := createHmac(queryString, masterCredentials.SecretKey)

// 	req, err := http.NewRequest("GET", url+"?"+queryString+"&signature="+signature, nil)
// 	if err != nil {
// 		log.Fatalf("Could not create request: %v", err)
// 	}
// 	req.Header.Add("X-MBX-APIKEY", masterCredentials.APIKey)

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Fatalf("Request failed: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Failed to read response body: %v", err)
// 	}

// 	// 응답 형식에 맞게 구조체 정의
// 	type AssetInfo struct {
// 		Asset                  string  `json:"asset"`
// 		WalletBalance          float64 `json:"walletBalance,string"`
// 		UnrealizedProfit       float64 `json:"unrealizedProfit,string"`
// 		MarginBalance          float64 `json:"marginBalance,string"`
// 		MaintenanceMargin      float64 `json:"maintenanceMargin,string"`
// 		InitialMargin          float64 `json:"initialMargin,string"`
// 		PositionInitialMargin  float64 `json:"positionInitialMargin,string"`
// 		OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"`
// 		MaxWithdrawAmount      float64 `json:"maxWithdrawAmount,string"`
// 	}
// 	type ResponseData struct {
// 		Assets []AssetInfo `json:"assets"`
// 	}
// 	var responseData ResponseData

// 	if err := json.Unmarshal(body, &responseData); err != nil {
// 		log.Fatalf("JSON unmarshalling failed: %s", err)
// 	}

// 	for _, asset := range responseData.Assets {
// 		if asset.Asset == "USDT" {
// 			return int(asset.WalletBalance), nil
// 		}
// 	}
// 	return 0, fmt.Errorf("USDT not found in the response")
// }
