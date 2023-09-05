package handlers

type PriceQuantity struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type TradingViewPayload struct {
	Account      string `json:"account"`
	Symbol       string `json:"symbol"`
	PositionSide string `json:"positionSide"`
	Amount       int    `json:"amount"`
	Entry        bool   `json:"entry,omitempty"`
}

type AllStopLossTakeProfitPayload struct {
	Account string   `json:"account"`
	Symbol  string   `json:"symbol"`
	TP      *float64 `json:"tp,omitempty"`
	SL      *float64 `json:"sl,omitempty"`
}

type PartialTakeProfitPayload struct {
	Account      string         `json:"account"`
	Symbol       string         `json:"symbol"`
	PositionSide string         `json:"positionSide"`
	TP           *PriceQuantity `json:"tp,omitempty"`
}

type CloseOrderPayload struct {
	Account string  `json:"account"`
	Symbol  string  `json:"symbol"`
	Close   float64 `json:"close"`
}

type CloseAllOrderPayload struct {
	Account string `json:"account"`
}

func NewTradingViewPayload() TradingViewPayload {
	return TradingViewPayload{
		Entry: true,
	}
}
