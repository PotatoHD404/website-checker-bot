package ssm

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"os"
)

func Init() {
	// create ssm options
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = os.Getenv("REGION")
		return nil
	})
	ssmClient := ssm.NewFromConfig(cfg)
	// get ssm parameter
	param, err := ssmClient.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           aws.String(os.Getenv("TOKEN_PARAMETER")),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		fmt.Println(err)
		panic("can't get ssm parameter")
	}
	err = os.Setenv("BOT_TOKEN", *param.Parameter.Value)
	if err != nil {
		fmt.Println(err)
		panic("can't set environment variable")
	}
}
