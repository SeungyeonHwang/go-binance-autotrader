package main

import (
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/binance"
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/webhook"
)

func main() {

	//TODO:
	// Lambda + Echo Proxy 통합 형태로 엔드포인트 관리하고 람다 1개로 온 리퀘스트를 각 엔드포인트에 맞게 분기시킨다

	//TODO:
	//자산 저장 기능(추이 tracking)
	//자동 매매 트레이딩 프로그램을 개발하고 있는데 파일 디렉토리에 파일DB 를 추가해서 자산의 추이를 기록하고 싶어 일단 DB를 만드는 것부터 시작해야되 프로그램이 시작할때 없으면 생성되어야 하고
	// 구체적으로 필요한 데이터는 아래와 같아

	// 추적하려는 계정의 자산은 4개야 master, sub1, sub2, sub3
	// 이것의 자산을 하루하루 기록해야되 그래서 하루전 일주일전 한달전 자산추이를 비교할 수 있어
	// 어떻게 설계하는게 좋을까?

	//TODO:
	//잔액 조회 기능(파일 DB로 자산 추이)
	binance.GetFuturesBalance("master")
	binance.GetFuturesBalance("sub1", "hwang.sy.test.1@gmail.com")
	binance.GetFuturesBalance("sub2", "hwang.sy.test.2@gmail.com")
	binance.GetFuturesBalance("sub3", "hwang.sy.test.3@gmail.com")

	//TODO:
	//DELETE DB

	//TODO:
	//포지션 조회 기능

	//TODO:
	//TP/SL 기능

	//TODO:
	//Market Clsoe 기능

	//TODO:
	//Market Entry 기능

	//TODO:
	//webHook 기능(Trading View webhook에 의해 트리거됨, order를 실행한다)(물타기)
	webhook.StartWebServer()
}
