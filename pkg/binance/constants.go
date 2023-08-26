package binance

// Position Info
const (
	MASTER_TIME     = "1D(+1H)"
	MASTER_LEVERAGE = 3
	MASTER_AMOUNT   = 0.2
	MASTER_METHOD   = "BTC:ETH"

	SUB1_TIME     = "1H(+5min)"
	SUB1_LEVERAGE = 10
	SUB1_AMOUNT   = 0.2
	SUB1_METHOD   = "SWING1"

	SUB2_TIME     = "4H(+15min)"
	SUB2_LEVERAGE = 5
	SUB2_AMOUNT   = 0.1
	SUB2_METHOD   = "SWING2"

	SUB3_TIME     = "4H(+15min)"
	SUB3_LEVERAGE = 5
	SUB3_AMOUNT   = 0.1
	SUB3_METHOD   = "SWING2"
)

const (
	baseURL      = "https://fapi.binance.com"
	SLACK_SUB1   = "https://hooks.slack.com/services/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"
	SLACK_SUB2   = "https://hooks.slack.com/services/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"
	SLACK_SUB3   = "https://hooks.slack.com/services/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"
	SLACK_MASTER = "https://hooks.slack.com/services/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"
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
