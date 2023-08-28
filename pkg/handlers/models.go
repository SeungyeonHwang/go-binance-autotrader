package handlers

type PriceQuantity struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type TradingViewPayload struct {
	Account      string `json:"account"`
	Symbol       string `json:"symbol"`
	PositionSide string `json:"positionSide"`
	Leverage     int    `json:"leverage"`
	Amount       int    `json:"amount"`
	Entry        bool   `json:"entry,omitempty"`
}

type AllStopLossTakeProfitPayload struct {
	Account      string  `json:"account"`
	Symbol       string  `json:"symbol"`
	PositionSide string  `json:"positionSide"`
	TP           float64 `json:"tp"`
	SL           float64 `json:"sl"`
}

type PartialTakeProfitPayload struct {
	Account      string         `json:"account"`
	Symbol       string         `json:"symbol"`
	PositionSide string         `json:"positionSide"`
	TP           *PriceQuantity `json:"tp,omitempty"`
}
