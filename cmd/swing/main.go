package main

import (
	"context"
	"log"

	"github.com/SeungyeonHwang/go-binance-autotrader/cmd/swing/router"
	"github.com/SeungyeonHwang/go-binance-autotrader/config"
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
)

var (
	echolambda *echoadapter.EchoLambda
	configData *config.Config
)

func init() {
	sess := session.Must(session.NewSession())
	ssmClient := ssm.New(sess)

	var err error
	configData, err = config.LoadConfigurationFromSSM(ssmClient)
	if err != nil {
		log.Fatalf("Failed to load config from SSM: %v", err)
	}

	e := echo.New()
	handler := &handlers.Handler{Echo: e, Config: configData}
	router.SetUp(e, handler)
	echolambda = echoadapter.New(e)
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echolambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(HandleRequest)
}
