package binance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2/futures"
)

func trimQuantity(quantity, stepSize float64) float64 {
	factor := math.Pow(10, float64(getPrecision(stepSize)))
	rounded := math.Round(quantity*factor) / factor
	return rounded
}

func GetStepSizeForSymbol(client *futures.Client, symbol string) (float64, error) {
	info, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return 0, err
	}
	for _, s := range info.Symbols {
		if s.Symbol == symbol {
			for _, f := range s.Filters {
				if f["filterType"] == "LOT_SIZE" {
					return strconv.ParseFloat(f["stepSize"].(string), 64)
				}
			}
		}
	}
	return 0, fmt.Errorf("step size for symbol %s not found", symbol)
}

func GetTickSizeForSymbol(client *futures.Client, symbol string) (float64, error) {
	info, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return 0, err
	}
	for _, s := range info.Symbols {
		if s.Symbol == symbol {
			for _, f := range s.Filters {
				if f["filterType"] == "PRICE_FILTER" {
					return strconv.ParseFloat(f["tickSize"].(string), 64)
				}
			}
		}
	}
	return 0, fmt.Errorf("tick size for symbol %s not found", symbol)
}

func getCurrentFuturesPrice(client *futures.Client, symbol string) (float64, error) {
	log.Printf("Fetching current price for symbol: %s", symbol)
	prices, err := client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return 0, err
	}
	for _, p := range prices {
		if p.Symbol == symbol {
			return strconv.ParseFloat(p.Price, 64)
		}
	}
	return 0, fmt.Errorf("price for symbol %s not found", symbol)
}

func FormatSymbol(s string) string {
	s = strings.ToUpper(s)
	if strings.Contains(s, "USDT") {
		return strings.Split(s, "USDT")[0] + "USDT"
	}
	return s + "USDT"
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

type SlackPayload struct {
	Text string `json:"text"`
}

func SendSlackNotification(webhookUrl, msg string) error {
	slackBody, _ := json.Marshal(SlackPayload{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return err
	}
	return nil
}

func getPrecision(value float64) int {
	decimalStr := fmt.Sprintf("%.10f", value)
	parts := strings.Split(decimalStr, ".")
	if len(parts) != 2 {
		return 0
	}

	decimalPart := parts[1]
	precision := 0
	for i := len(decimalPart) - 1; i >= 0; i-- {
		if decimalPart[i] != '0' {
			precision = i + 1
			break
		}
	}

	return precision
}
