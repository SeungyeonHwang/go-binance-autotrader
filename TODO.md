TODO:
	[출력물 개선]
	- position / balance / order 출력 개선
	- Balance 줄 맞추기 / 표시 좀 다듬기
	- 주문 정보 USDT 정보 및 필요 정보 표시하기
=========================================================================================================================================================
지금 당장 필요한 기능만 개발 
나머진 프론트 개발하면서
스탑 티피 기능 즉시 엔트리 기능
핸드폰에 슬랙 넣기
dev -> pr 승인되면 자동 배포
TODO: 
	- 잔액 파악해서 처음 UnitPrice 설정 -> Unitprice 지정할 필요 없어진.
	- 현재 포지션 정보 들고와서 +면 그대로 / -면 반만 물타게 변경

TODO:
	- 람다에서 파일 DB 가능 한가 ?
TODO:
	잔액 확인, 주문 확인 Beautify

TODO:
	Unit Price 추적 및 추이 알 수 있게 하기

TODO:
	포지션 정보 가져오기 GET

		 TODO
amount -> unitprice 전체 자산의 5% *0.05 X Leverage
파일 DB에 저장해놨다가 거래 시작할때 읽어 놓고 
+ 면 1/2 -면 1

	TODO:
		잔액 표시 기능

 		-> 자산 추이를 보여주게 하고싶다
		 -> 다른 정보 필요없음 추이만 알면 됨


TODO
clear unitprice 기능
	//현재의 잔액을 계산
	//UnitPrice 도출한다 전체 자산의 5% *0.05 X Leverage
	//파일 DB에 저장 한다.
=========================================================================================================================================================
	//TODO:
		Cognito로 로그인 한사람만 가능 하게 제안
		-> 로그인 한사람의 secret 랑 연결

	//TODO:
		테스트 코드 추가(정상계)

	//TODO:
		슬랙 웹 훅 나누기

	//TODO: 
	포지션 정보

	//TODO:
	//자산 저장 기능(추이 tracking)(파일 DB로 자산 추이)
	//자동 매매 트레이딩 프로그램을 개발하고 있는데 파일 디렉토리에 파일DB 를 추가해서 자산의 추이를 기록하고 싶어 일단 DB를 만드는 것부터 시작해야되 프로그램이 시작할때 없으면 생성되어야 하고
	// 구체적으로 필요한 데이터는 아래와 같아

	// 추적하려는 계정의 자산은 4개야 master, sub1, sub2, sub3
	// 이것의 자산을 하루하루 기록해야되 그래서 하루전 일주일전 한달전 자산추이를 비교할 수 있어
	// 어떻게 설계하는게 좋을까?

	//TODO:
	//DELETE DB

	//TODO:
	//TP/SL 기능

	//TODO:
	//Market Clsoe 기능

	//TODO:
	//Market Entry 기능

	//TODO:
	//webHook 기능(Trading View webhook에 의해 트리거됨, order를 실행한다)(물타기)
	//webhook.StartWebServer()

	//TODO:
	리액트로 UI 배포하기
	도메인 사기

	//TODO:
	Cognito 로 로그인 연결 하고 
	apikey, secretkey 관리 DB로 하기

	//TODO:
	Dev 환경 만들기

	트레이딩 봇 기능에 집중하기(다른기능 넣어봣자 바이낸스 본사이트 이용하면 되기 떄문에)
DB도입
Order ID로 ST/TP 관리 하기 , 부분 익절 같은거



코인판에 판매 하기
포트폴리오로 기록하기(서류 작성하기)

=========================================================================================================================================================

# Pumping
현물 (Pumping) : upbit, bithumb 추가

# Grid(Auto)
선물 : Grid 자산 조회 기능
    - bitget
    - bybit
    - okx

# Defi

# 차익거래