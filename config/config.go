package config

import (
	"github.com/SeungyeonHwang/go-binance-autotrader/pkg/handlers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMConfigLoader struct {
	SSMClient *ssm.SSM
}

func (loader *SSMConfigLoader) GetParameter(paramName string) (string, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	}
	result, err := loader.SSMClient.GetParameter(input)
	if err != nil {
		return "", err
	}

	return *result.Parameter.Value, nil
}

func LoadConfigurationFromSSM(ssmClient *ssm.SSM) (*handlers.Config, error) {
	loader := SSMConfigLoader{SSMClient: ssmClient}
	return loadConfigFromSSM(loader)
}

func loadConfigFromSSM(loader SSMConfigLoader) (*handlers.Config, error) {
	config := &handlers.Config{}

	var err error
	config.MasterAPIKey, err = loader.GetParameter("/binance/master/api_key")
	if err != nil {
		return nil, err
	}

	config.MasterSecretKey, err = loader.GetParameter("/binance/master/secret_key")
	if err != nil {
		return nil, err
	}

	config.Sub1APIKey, err = loader.GetParameter("/binance/sub1/api_key")
	if err != nil {
		return nil, err
	}

	config.Sub1SecretKey, err = loader.GetParameter("/binance/sub1/secret_key")
	if err != nil {
		return nil, err
	}

	return config, nil
}
