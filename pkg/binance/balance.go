package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/SeungyeonHwang/go-binance-autotrader/config"
)

func FetchAllBalances(config *config.Config) (string, error) {
	var resultBuilder strings.Builder

	for _, acc := range Accounts {
		balance, err := GetFuturesBalance(acc.AccountType, config, acc.Email)
		if err != nil {
			return "", err
		}

		resultBuilder.WriteString(":bank: " + acc.Label + "\n")
		resultBuilder.WriteString(strings.Repeat("-", 40) + "\n")

		if info, ok := SubInfos[acc.Label]; ok {
			unitPrice := int(float64(balance) * info.Amount * float64(info.Leverage))
			resultBuilder.WriteString(":clock1: Time: " + info.Time + "\n")
			resultBuilder.WriteString(":hammer_and_pick: Method: " + info.Method + "\n")
			resultBuilder.WriteString(":rocket: Leverage: x" + strconv.Itoa(info.Leverage) + "\n")
			resultBuilder.WriteString(":dollar: Unit Price: $" + fmt.Sprintf("%d", unitPrice) + "\n")
		}

		resultBuilder.WriteString(":moneybag: Balance: $" + fmt.Sprintf("%d", balance) + "\n")
		resultBuilder.WriteString(strings.Repeat("-", 40) + "\n")
		resultBuilder.WriteString("\n\n")
	}

	return resultBuilder.String(), nil
}

func GetFuturesBalance(account string, config *config.Config, subAccountEmail ...string) (int, error) {
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
		apiKey      string
		secretKey   string
	}{
		{MASTER_ACCOUNT, "Master", config.MasterAPIKey, config.MasterSecretKey},
		{SUB1_ACCOUNT, "Sub1", config.Sub1APIKey, config.Sub1SecretKey},
		{SUB2_ACCOUNT, "Sub2", config.Sub2APIKey, config.Sub2SecretKey},
		{SUB3_ACCOUNT, "Sub3", config.Sub3APIKey, config.Sub3SecretKey},
	}

	var resultBuilder strings.Builder
	lineSeparator := strings.Repeat("-", 40) + "\n"

	for _, acc := range accounts {
		positions, totalCrossUnPnl, availableBalance, totalInitialMargin, err := fetchPositions(acc.apiKey, acc.secretKey)
		if err != nil {
			return "", err
		}

		resultBuilder.WriteString(":bank: " + acc.label + "\n")
		resultBuilder.WriteString(lineSeparator)

		sign := ""
		if totalCrossUnPnl > 0 {
			sign = "+"
		}
		resultBuilder.WriteString("[Profit]: " + sign + fmt.Sprintf("%.1f", totalCrossUnPnl) + "\n")

		if totalInitialMargin == 0 {
			resultBuilder.WriteString("[ROI]: 0\n")
		} else {
			roi := (totalCrossUnPnl / totalInitialMargin) * 100

			sign := ""
			if roi > 0 {
				sign = "+"
			}
			resultBuilder.WriteString("[ROI]: " + sign + fmt.Sprintf("%.1f", roi) + "%\n")
		}

		resultBuilder.WriteString("[Available]: " + fmt.Sprintf("%.1f", availableBalance) + "\n")
		resultBuilder.WriteString(lineSeparator)

		for _, position := range positions {
			amt, err := strconv.ParseFloat(position.PositionAmt, 64)
			if err != nil || amt == 0 {
				continue
			}

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
			resultBuilder.WriteString(position.Symbol + ": " + profitStr + "\n")
		}

		resultBuilder.WriteString(lineSeparator)
		resultBuilder.WriteString("\n")
	}

	return resultBuilder.String(), nil
}

func fetchPositions(apiKey string, secretKey string) ([]Position, float64, float64, float64, error) {
	url := baseURL + "/fapi/v2/account"

	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	queryString := "timestamp=" + timestamp
	signature := createHmac(queryString, secretKey)

	fullURL := url + "?" + queryString + "&signature=" + signature
	log.Println("Requesting:", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, 0, 0, 0, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Add("X-MBX-APIKEY", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return nil, 0, 0, 0, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, 0, 0, 0, fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData struct {
		TotalCrossUnPnl    string     `json:"totalCrossUnPnl"`
		AvailableBalance   string     `json:"availableBalance"`
		TotalInitialMargin string     `json:"totalInitialMargin"`
		Positions          []Position `json:"positions"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return nil, 0, 0, 0, err
	}

	totalCrossUnPnl, err := strconv.ParseFloat(responseData.TotalCrossUnPnl, 64)
	if err != nil {
		log.Printf("Error converting totalCrossUnPnl to float: %v", err)
		return nil, 0, 0, 0, fmt.Errorf("failed to convert totalCrossUnPnl: %v", err)
	}

	availableBalance, err := strconv.ParseFloat(responseData.AvailableBalance, 64)
	if err != nil {
		log.Printf("Error converting availableBalance to float: %v", err)
		return nil, 0, 0, 0, fmt.Errorf("failed to convert availableBalance: %v", err)
	}

	totalInitialMargin, err := strconv.ParseFloat(responseData.TotalInitialMargin, 64)
	if err != nil {
		log.Printf("Error converting totalInitialMargin to float: %v", err)
		return nil, 0, 0, 0, fmt.Errorf("failed to convert totalInitialMargin: %v", err)
	}

	sort.Slice(responseData.Positions, func(i, j int) bool {
		profitI, errI := strconv.ParseFloat(responseData.Positions[i].UnrealizedProfit, 64)
		profitJ, errJ := strconv.ParseFloat(responseData.Positions[j].UnrealizedProfit, 64)

		if errI != nil || errJ != nil {
			return false
		}
		return profitI > profitJ
	})

	return responseData.Positions, totalCrossUnPnl, availableBalance, totalInitialMargin, nil
}

func positionExistsForSymbol(apiKey, secretKey, targetSymbol string) (bool, error) {
	positions, _, _, _, err := fetchPositions(apiKey, secretKey)
	if err != nil {
		return false, err
	}

	for _, position := range positions {
		if strings.EqualFold(position.Symbol, targetSymbol) {
			amt, err := strconv.ParseFloat(position.PositionAmt, 64)
			if err != nil {
				continue
			}
			if amt != 0 {
				return true, nil
			}
		}
	}
	return false, nil
}

func getROIForSymbol(apiKey, secretKey, targetSymbol string) (float64, error) {
	exists, err := positionExistsForSymbol(apiKey, secretKey, targetSymbol)
	if err != nil || !exists {
		return 0, fmt.Errorf("position for symbol %s not found", targetSymbol)
	}

	positions, _, _, _, err := fetchPositions(apiKey, secretKey)
	if err != nil {
		return 0, err
	}

	for _, position := range positions {
		if strings.EqualFold(position.Symbol, targetSymbol) {
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
	return 0, fmt.Errorf("position for symbol %s not found", targetSymbol)
}
