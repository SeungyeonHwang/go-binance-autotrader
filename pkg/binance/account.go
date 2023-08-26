package binance

type AccountConfig struct {
	APIKey    string
	SecretKey string
}

type Account struct {
	AccountType string
	Email       string
	Label       string
}

type SubInfo struct {
	Amount   float64
	Leverage int
	Time     string
	Method   string
}

var Accounts = []Account{
	{MASTER_ACCOUNT, "", "Master"},
	{SUB1_ACCOUNT, SUB1_EMAIL, "Sub1"},
	{SUB2_ACCOUNT, SUB2_EMAIL, "Sub2"},
	{SUB3_ACCOUNT, SUB3_EMAIL, "Sub3"},
}

var SubInfos = map[string]SubInfo{
	"Master": {MASTER_AMOUNT, MASTER_LEVERAGE, MASTER_TIME, MASTER_METHOD},
	"Sub1":   {SUB1_AMOUNT, SUB1_LEVERAGE, SUB1_TIME, SUB1_METHOD},
	"Sub2":   {SUB2_AMOUNT, SUB2_LEVERAGE, SUB2_TIME, SUB2_METHOD},
	"Sub3":   {SUB3_AMOUNT, SUB3_LEVERAGE, SUB3_TIME, SUB3_METHOD},
}
