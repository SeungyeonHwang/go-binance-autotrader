package binance

// Position Info
const (
	MASTER_TIME     = "1D"
	MASTER_LEVERAGE = 10
	MASTER_AMOUNT   = 0.25
	MASTER_METHOD   = "BTC/ETH"

	SUB1_TIME     = "1D"
	SUB1_LEVERAGE = 10
	SUB1_AMOUNT   = 0.3
	SUB1_METHOD   = "ALT"

	SUB2_TIME     = "1D"
	SUB2_LEVERAGE = 10
	SUB2_AMOUNT   = 0.3
	SUB2_METHOD   = "ALT"

	SUB3_TIME     = "1D"
	SUB3_LEVERAGE = 10
	SUB3_AMOUNT   = 0.3
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
