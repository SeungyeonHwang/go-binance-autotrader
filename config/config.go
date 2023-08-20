package config

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Config struct {
	Binance map[string]struct {
		APIKey    string
		SecretKey string
	}
}

func GetConfig(sess *session.Session) (Config, error) {
	ssmClient := ssm.New(sess)

	paramNames := []string{
		"/binance/master/api_key",
		"/binance/master/secret_key",
		"/binance/sub1/api_key",
		"/binance/sub1/secret_key",
	}

	paramInput := &ssm.GetParametersInput{
		Names:          aws.StringSlice(paramNames),
		WithDecryption: aws.Bool(true),
	}

	resp, err := ssmClient.GetParameters(paramInput)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		Binance: make(map[string]struct {
			APIKey    string
			SecretKey string
		}),
	}

	for _, param := range resp.Parameters {
		parts := strings.Split(aws.StringValue(param.Name), "/")
		if len(parts) < 4 {
			continue
		}
		accountType := parts[2]
		keyType := parts[3]

		acc := config.Binance[accountType]

		if keyType == "api_key" {
			acc.APIKey = aws.StringValue(param.Value)
		} else if keyType == "secret_key" {
			acc.SecretKey = aws.StringValue(param.Value)
		}

		config.Binance[accountType] = acc
	}

	return config, nil
}
