.PHONY: build

AWS_PROFILE=hwang_personal

build:
	@echo "Building and Deploying the BinanceTradingService..."
	sam build
	sam deploy --profile $(AWS_PROFILE)
