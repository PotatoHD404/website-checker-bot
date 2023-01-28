package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/tucnak/telebot.v2"
	"net/http"
	"os"
)

var err error = nil

var tgBot *telebot.Bot

func initTelebot() {
	settings := telebot.Settings{
		Token:       os.Getenv("BOT_TOKEN"),
		Synchronous: true,
		Verbose:     true,
	}
	tgBot, err = telebot.NewBot(settings)
	if err != nil {
		fmt.Println(err)
		panic("can't create bot")
	}
	tgBot.Handle(telebot.OnText, func(m *telebot.Message) {
		message := m.Text
		_, err := tgBot.Send(m.Sender, message)
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})
}

var svc *dynamodb.DynamoDB

func initDynamodb() {
	dynamodbSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc = dynamodb.New(dynamodbSession)
}

var adapter *httpadapter.HandlerAdapter

func initRouter() {
	r := httprouter.New()
	r.HandlerFunc("GET", "/setWebhook", setWebhook)
	r.HandlerFunc("POST", "/handleMessage", handleMessage)
	adapter = httpadapter.New(r)
}

func init() {
	initTelebot()
	initDynamodb()
	initRouter()
}

func setWebhook(w http.ResponseWriter, r *http.Request) {
	// Set webhook
	err := tgBot.SetWebhook(&telebot.Webhook{
		Listen: "",
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: "https://api.telegram.org/bot" + os.Getenv("BOT_TOKEN") + "/",
		},
	})
	if err != nil {
		fmt.Println(err)
		panic("can't set webhook")
	}

	returnOk(w)
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	// Process update
	var u telebot.Update
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println(err)
		panic("can't unmarshal")
	}
	tgBot.ProcessUpdate(u)
	//	return body: "ok", statusCode: 200
	returnOk(w)
}

func returnOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		fmt.Println(err)
		panic("can't write response")
	}
}

func main() {
	lambda.Start(adapter.ProxyWithContext)
}
