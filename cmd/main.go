package main

import "github.com/SeungyeonHwang/go-binance-autotrader/pkg/binance"

func main() {

	binance.GetFuturesBalance("master")
	binance.GetFuturesBalance("sub1", "hwang.sy.test.1@gmail.com")
	// binance.GetFuturesBalance("sub2", "hwang.sy.test.2@gmail.com")
	// binance.GetFuturesBalance("sub3", "hwang.sy.test.3@gmail.com")
	// webhook.StartWebServer()
}
