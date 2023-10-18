package binance

// Position Info
const (
	TIME         = "4H"
	LEVERAGE     = 7
	LEVERAGE_ALT = 5

	MASTER_TIME     = TIME
	MASTER_LEVERAGE = LEVERAGE
	MASTER_AMOUNT   = 0.25
	MASTER_METHOD   = "BTC/ETH"

	SUB1_TIME     = TIME
	SUB1_LEVERAGE = LEVERAGE_ALT
	SUB1_AMOUNT   = 0.4
	SUB1_METHOD   = "ALT"

	SUB2_TIME     = TIME
	SUB2_LEVERAGE = LEVERAGE_ALT
	SUB2_AMOUNT   = 0.4
	SUB2_METHOD   = "ALT"

	SUB3_TIME     = TIME
	SUB3_LEVERAGE = LEVERAGE_ALT
	SUB3_AMOUNT   = 0.4
	SUB3_METHOD   = "ALT"
)

const (
	baseURL = "https://fapi.binance.com"
)

const (
	BUCKET_NAME = "asset-balance-bucket"
	DB_NAME     = "asset-database"
)

const (
	MASTER_ACCOUNT = "master"
	SUB1_ACCOUNT   = "sub1"
	SUB1_EMAIL     = "hwang.sy.test.1@gmail.com"
	SUB2_ACCOUNT   = "sub2"
	SUB2_EMAIL     = "hwang.sy.test.2@gmail.com"
	SUB3_ACCOUNT   = "sub3"
	SUB3_EMAIL     = "hwang.sy.test.3@gmail.com"
)

const (
	OPEN  = "OPEN"
	CLOSE = "CLOSE"

	LONG  = "LONG"
	SHORT = "SHORT"

	CROSSED = "CROSSED"
)
