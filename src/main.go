package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"website-checker-bot/router"
)

var adapter *httpadapter.HandlerAdapterV2

func init() {
	adapter = router.GetAdapter()
}

func main() {
	lambda.Start(adapter.ProxyWithContext)
}
