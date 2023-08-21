package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SeungyeonHwang/go-binance-autotrader/config"
)

type PositionInfo struct {
	Asset            string `json:"asset"`
	UnrealizedProfit string `json:"unrealizedProfit"`
}

func FetchAllBalances(config *config.Config) (string, error) {
	accounts := []struct {
		accountType string
		email       string
		label       string
	}{
		{MASTER_ACCOUNT, "", "Master"},
		{SUB1_ACCOUNT, SUB1_EMAIL, "Sub1"},
		// {MASTER_ACCOUNT, SUB2_EMAIL, "Sub2"},
		// {MASTER_ACCOUNT, SUB3_EMAIL, "Sub3"},
	}

	var resultBuilder strings.Builder

	for _, acc := range accounts {
		balance, err := getFuturesBalance(acc.accountType, config, acc.email)
		if err != nil {
			return "", err
		}

		resultBuilder.WriteString(fmt.Sprintf("Account     | %s\n", acc.label))
		resultBuilder.WriteString("--------------------\n")

		if acc.label == "Sub1" {
			unitPrice := int(float64(balance) * 0.05 * 15)
			resultBuilder.WriteString(fmt.Sprintf("Time        | 1H\n"))
			resultBuilder.WriteString(fmt.Sprintf("Leverage    | X15\n"))
			resultBuilder.WriteString(fmt.Sprintf("Unit Price  | %d\n", unitPrice))
			resultBuilder.WriteString(fmt.Sprintf("Balance     | *%d\n", balance))
		} else {
			resultBuilder.WriteString(fmt.Sprintf("Balance     | *%d\n", balance))
		}
		resultBuilder.WriteString("--------------------\n")
		resultBuilder.WriteString("\n")
	}

	return resultBuilder.String(), nil
}

func FetchAllPositions(config *config.Config) (string, error) {
	accounts := []struct {
		accountType string
		email       string
		label       string
	}{
		// {MASTER_ACCOUNT, "", "Master"},
		{SUB1_ACCOUNT, SUB1_EMAIL, "Sub1"},
		// {MASTER_ACCOUNT, SUB2_EMAIL, "Sub2"},
		// {MASTER_ACCOUNT, SUB3_EMAIL, "Sub3"},
	}

	var resultBuilder strings.Builder

	for _, acc := range accounts {
		positions, err := fetchPositions(config.MasterAPIKey, config.MasterSecretKey, acc.email)
		if err != nil {
			return "", err
		}

		resultBuilder.WriteString(fmt.Sprintf("%-15s| %s\n", "Account", acc.label))
		resultBuilder.WriteString("--------------------\n")

		for _, position := range positions {
			profit, _ := strconv.ParseFloat(position.UnrealizedProfit, 64)
			resultBuilder.WriteString(fmt.Sprintf("%-15s| Unrealized Profit: %.1f\n", "Asset: "+position.Asset, profit))
		}

		resultBuilder.WriteString("--------------------\n")
		resultBuilder.WriteString("\n")
	}

	return resultBuilder.String(), nil
}

func getFuturesBalance(account string, config *config.Config, subAccountEmail ...string) (int, error) {
	if account == MASTER_ACCOUNT {
		return fetchBalance(config.MasterAPIKey, config.MasterSecretKey)
	} else {
		return fetchSubAccountBalance(config.MasterAPIKey, config.MasterSecretKey, subAccountEmail[0])
	}
}

func fetchBalance(apiKey string, secretKey string) (int, error) {
	url := baseURL + "/fapi/v2/balance"
	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	queryString := "timestamp=" + timestamp
	signature := createHmac(queryString, secretKey)

	fullURL := url + "?" + queryString + "&signature=" + signature
	log.Println("Requesting:", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
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

func fetchSubAccountBalance(apiKey string, secretKey string, subAccountEmail string) (int, error) {
	url := "https://api.binance.com/sapi/v1/sub-account/futures/account"
	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	queryString := "timestamp=" + timestamp + "&email=" + subAccountEmail
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

	type AssetInfo struct {
		Asset         string  `json:"asset"`
		WalletBalance float64 `json:"walletBalance,string"`
	}
	type ResponseData struct {
		Assets []AssetInfo `json:"assets"`
	}
	var responseData ResponseData

	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("JSON unmarshalling error. Body: %s, Error: %v", string(body), err)
		return 0, fmt.Errorf("JSON unmarshalling failed: %s", err)
	}

	for _, asset := range responseData.Assets {
		if asset.Asset == "USDT" {
			return int(asset.WalletBalance), nil
		}
	}
	log.Println("USDT not found in the response")
	return 0, fmt.Errorf("USDT not found in the response")
}

func fetchPositions(apiKey string, secretKey string, subAccountEmail string) ([]PositionInfo, error) {
	var url string
	if subAccountEmail == "" {
		url = baseURL + "/fapi/v2/account"
	} else {
		url = "https://api.binance.com/sapi/v1/sub-account/futures/account"
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	queryString := "timestamp=" + timestamp
	if subAccountEmail != "" {
		queryString += "&email=" + subAccountEmail
	}
	signature := createHmac(queryString, secretKey)

	fullURL := url + "?" + queryString + "&signature=" + signature
	log.Println("Requesting:", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Add("X-MBX-APIKEY", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData struct {
		Assets []PositionInfo `json:"assets"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("JSON unmarshalling error. Body: %s, Error: %v", string(body), err)
		return nil, fmt.Errorf("JSON unmarshalling failed: %s", err)
	}

	var assetsInfo []PositionInfo
	for _, position := range responseData.Assets {
		assetsInfo = append(assetsInfo, PositionInfo{
			Asset:            position.Asset,
			UnrealizedProfit: position.UnrealizedProfit,
		})
	}

	return assetsInfo, nil
}
