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

type Position struct {
	Symbol           string `json:"symbol"`
	InitialMargin    string `json:"initialMargin"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	PositionAmt      string `json:"positionAmt"`
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

func FetchAllPositions(config *config.Config) (string, error) {
	accounts := []struct {
		accountType string
		label       string
	}{
		// {MASTER_ACCOUNT, "Master"},
		{SUB1_ACCOUNT, "Sub1"},
		// {MASTER_ACCOUNT, "Sub2"},
		// {MASTER_ACCOUNT, "Sub3"},
	}

	var resultBuilder strings.Builder

	for _, acc := range accounts {
		var apiKey, secretKey string
		switch acc.accountType {
		case MASTER_ACCOUNT:
			apiKey = config.MasterAPIKey
			secretKey = config.MasterSecretKey
		case SUB1_ACCOUNT:
			apiKey = config.Sub1APIKey
			secretKey = config.Sub1SecretKey
		default:
			return "", fmt.Errorf("unknown account type: %s", acc.accountType)
		}

		positions, totalCrossUnPnl, err := fetchPositions(apiKey, secretKey)
		if err != nil {
			return "", err
		}

		resultBuilder.WriteString(fmt.Sprintf("%-15s | Total: %-20.1f\n", acc.label, totalCrossUnPnl))
		resultBuilder.WriteString("------------------------------------\n")

		for _, position := range positions {
			amt, err := strconv.ParseFloat(position.PositionAmt, 64)
			if err != nil {
				continue
			}
			if amt != 0 {
				profit, errProfit := strconv.ParseFloat(position.UnrealizedProfit, 64)
				initialMargin, errMargin := strconv.ParseFloat(position.InitialMargin, 64)

				if errProfit != nil || errMargin != nil {
					continue
				}

				roi := 0.0
				if initialMargin != 0 {
					roi = (profit / initialMargin) * 100
				}

				profitStr := fmt.Sprintf("%.1f (%.2f%%)", profit, roi)
				if profit > 0 {
					profitStr = fmt.Sprintf("+%.1f (+%.2f%%)", profit, roi)
				}
				resultBuilder.WriteString(fmt.Sprintf("%-15s | %-30s\n", position.Symbol, profitStr))
			}
		}

		resultBuilder.WriteString("------------------------------------\n")
		resultBuilder.WriteString("\n")

	}
	return resultBuilder.String(), nil
}

func fetchPositions(apiKey string, secretKey string) ([]Position, float64, error) {
	url := baseURL + "/fapi/v2/account"

	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	queryString := "timestamp=" + timestamp
	signature := createHmac(queryString, secretKey)

	fullURL := url + "?" + queryString + "&signature=" + signature
	log.Println("Requesting:", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, 0, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Add("X-MBX-APIKEY", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return nil, 0, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, 0, fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData struct {
		TotalCrossUnPnl string     `json:"totalCrossUnPnl"`
		Positions       []Position `json:"positions"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("Error unmarshalling response body: %v", err) // It's good to log the specific error.
		return nil, 0, err
	}

	totalCrossUnPnl, err := strconv.ParseFloat(responseData.TotalCrossUnPnl, 64)
	if err != nil {
		log.Printf("Error converting totalCrossUnPnl to float: %v", err)
		return nil, 0, fmt.Errorf("failed to convert totalCrossUnPnl: %v", err)
	}

	return responseData.Positions, totalCrossUnPnl, nil
}

func getROIForSymbol(apiKey, secretKey, targetSymbol string) (float64, error) {
	positions, _, err := fetchPositions(apiKey, secretKey)
	if err != nil {
		return 0, err
	}

	for _, position := range positions {
		if strings.EqualFold(position.Symbol, targetSymbol) {
			amt, err := strconv.ParseFloat(position.PositionAmt, 64)
			if err != nil {
				continue
			}
			if amt != 0 {
				profit, errProfit := strconv.ParseFloat(position.UnrealizedProfit, 64)
				initialMargin, errMargin := strconv.ParseFloat(position.InitialMargin, 64)

				if errProfit != nil || errMargin != nil {
					continue
				}

				if initialMargin == 0 {
					return 0, fmt.Errorf("initial margin for symbol %s is zero", targetSymbol)
				}

				return (profit / initialMargin) * 100, nil
			}
		}
	}
	return -1, fmt.Errorf("position for symbol %s not found", targetSymbol)
}
