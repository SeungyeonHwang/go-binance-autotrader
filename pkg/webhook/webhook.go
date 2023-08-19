package webhook

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/binance"
)

//	{
//	  "account":"sub1"
//	  "symbol": "BTCUSDT",
//	  "amount": 30.0,
//	  "positionSide": "SHORT",
//	}
type TradingViewPayload struct {
	Account      string  `json:"account"`
	Symbol       string  `json:"symbol"`
	Amount       float64 `json:"amount"`
	PositionSide string  `json:"positionSide"`
}

func StartWebServer() {
	http.HandleFunc("/getbalance", webhookHandler)
	http.HandleFunc("/webhook", webhookHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// POST 요청만 허용합니다.
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	// 요청 본문을 읽습니다.
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	var alertData TradingViewPayload

	err = json.Unmarshal(body, &alertData)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		return
	}

	// Logging parsed data
	log.Printf("Parsed Alert Data: Account: %s, Symbol: %s, Amount: %f, PositionSide: %s\n", alertData.Account, alertData.Symbol, alertData.Amount, alertData.PositionSide)

	// Binance 주문 실행
	err = binance.PlaceFuturesMarketOrder(alertData.Account, alertData.Symbol, alertData.Amount, alertData.PositionSide)
	if err != nil {
		log.Printf("Failed to place futures market order: %s", err)
		http.Error(w, "Error executing order", http.StatusInternalServerError)
		return
	}
}
