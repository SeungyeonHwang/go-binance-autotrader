package main

import (
	"context"

	"github.com/SeungyeonHwang/go-binance-autotrader/cmd/swing/router"
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
)

var (
	echolambda *echoadapter.EchoLambda
)

func init() {
	e := echo.New()
	handler := &handlers.Handler{Echo: e}
	router.SetUp(e, handler)
	echolambda = echoadapter.New(e)
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echolambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(HandleRequest)
}
