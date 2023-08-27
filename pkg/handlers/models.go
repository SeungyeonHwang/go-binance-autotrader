package handlers

type TradingViewPayload struct {
	Account      string `json:"account"`
	Symbol       string `json:"symbol"`
	PositionSide string `json:"positionSide"`
	Leverage     int    `json:"leverage"`
	Amount       int    `json:"amount"`
	Entry        bool   `json:"entry,omitempty"`
}

type StopLossTakeProfitPayload struct {
	Account      string  `json:"account"`
	Symbol       string  `json:"symbol"`
	PositionSide string  `json:"positionSide"`
	TP           float64 `json:"tp"`
	SL           float64 `json:"sl"`
}
