.PHONY: build

AWS_PROFILE=hwang_personal

build:
	@echo "Removing the .aws-sam directory..."
	rm -rf .aws-sam

	@echo "Building and Deploying the BinanceTradingService..."
	sam build
	sam deploy --profile $(AWS_PROFILE)
