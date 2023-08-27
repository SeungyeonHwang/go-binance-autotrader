// 2. Stop loss & Take Profit // Partial




		// 3. clear open order
		// 3. Trailing Stop
		// 4. Closs Position All(specific coin) 
			// clear open orders
		// 6. Closs Position All(all coin)
			// clear open orders	
			바이낸스 API 를 통해 거래하도록 프로그램을 설계하고있어 
구현하고 싶은 기능은 아래와같아 아래와같은 Json 을 payload 에 실어서 실행시키는 느낌으로
어떤식으로 API 를 구현하는게 좋을까 ? 단순 플래그나 Json 데이터를 조정하는 것으로 구현 API 를 줄일 수 있으면 그쪽 방향도 검토하고 있어 
일단 설계를 부탁할게
{
	"account": "sub1",
	"symbol": "XRPUSDT",
	....
}

1. Stop loss & Take Profit (Market)
2. Stop loss & Take Profit (Limit) 
3. Trailing Stop
4. Closs Position All(specific coin)
5. Closs Position Partial(specific coin)
6. Closs Position All(all coin)

Stop loss & Take Profit (Market)

Stop loss와 Take profit의 경우 주문 (Order)를 생성할 때 해당 옵션들을 지정할 수 있습니다.
POST /api/v3/order 를 사용하고 필요한 매개변수로 stopPrice와 price를 포함시킵니다.
Stop loss & Take Profit (Limit)

마찬가지로 주문 생성 API를 사용하지만 이번에는 timeInForce 매개변수를 GTC (Good Till Cancelled) 값으로 지정하여 제한 주문 (Limit order)으로 생성합니다.
Trailing Stop



GET /fapi/v1/openOrders: 모든 오픈 주문 조회

POST /fapi/v1/order: 주문 생성
POST /fapi/v1/batchOrders: 여러 주문 동시 생성
GET /fapi/v1/order: 주문 정보 조회
DELETE /fapi/v1/order: 주문 취소
DELETE /fapi/v1/batchOrders: 여러 주문 동시 취소
GET /fapi/v1/openOrder: 개별 오픈 주문 조회

GET /fapi/v1/allOrders: 모든 주문 정보 조회

바이낸스에는 직접적인 Trailing Stop 기능은 없습니다. 이 기능을 구현하려면 코드로 로직을 작성해야 합니다.
가격이 특정 퍼센트만큼 상승하면 Stop loss를 동일한 퍼센트만큼 올리는 로직을 구현해야 합니다.
Close Position All (specific coin)

특정 코인에 대한 모든 포지션을 닫으려면 해당 코인에 대한 현재 주문들을 확인 (GET /api/v3/openOrders)하고, 모든 주문을 취소 (DELETE /api/v3/order)합니다.
Close Position Partial (specific coin)

특정 코인에 대한 포지션을 부분적으로 닫으려면 주문 크기를 조절하여 주문을 생성하거나 변경해야 합니다.
Close Position All (all coin)

모든 코인에 대한 포지션을 닫으려면 모든 코인의 현재 주문들을 확인하고, 모든 주문을 취소해야 합니다. 이를 위해서는 모든 코인에 대해 주문 조회 및 주문 취소 API를 반복적으로 호출해야 합니다.

==================================================================================================================
- Close all / 비중 close
- ST/TP 기능
- Market Entry
- 레버리지 넘으면 안되는거 고치기
- dev -> pr 승인되면 자동 배포
- dev 환경만들기 (출시 예정) : 챗봇 기능에만 집중하기
- 업비트 / 빗썸 연동
- Cognito로 로그인 한사람만 가능 하게 제안
		-> 로그인 한사람의 secret 랑 연결
- 테스트 코드 추가(정상계)
- React UI 배포 / 도메인 사기
- 구글 광고
- 리퍼럴 넣기
- apikey, secretkey 관리 DB로 하기
- 트레이딩 봇 기능에 집중하기(다른기능 넣어봣자 바이낸스 본사이트 이용하면 되기 떄문에)
- Order ID로 ST/TP 관리 하기 , 부분 익절 같은거
- 광고 제거 등 플랜
- 거래소 늘리기, 서브어카운트 사용 기능 추가
- 코인판에 판매 하기
- 포트폴리오로 기록하기(서류 작성하기)
- 신기능 Dev에서 개발
- 엔트리 기능 넣기