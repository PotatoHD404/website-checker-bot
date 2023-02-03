package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"sync"
	"website-checker-bot/src/adapter"
	"website-checker-bot/src/dynamodb"
	"website-checker-bot/src/ssm"
	"website-checker-bot/src/telebot"
)

var Wg sync.WaitGroup

func init() {
	ssm.initSSM()

	Wg = sync.WaitGroup{}
	Wg.Add(3)

	go telebot.initTelebot()
	go dynamodb.initDynamodb()
	go adapter.initRouter()

	Wg.Wait()
}

func main() {
	lambda.Start(adapter.adapter.ProxyWithContext)
}
