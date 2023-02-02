package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	adminsTable      = "checker-admins"
	subscribersTable = "checker-subscribers"
)

func initSSM() {
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
		panic("can't set env variable")
	}
}

var tgBot *telebot.Bot

func initTelebot() {
	defer wg.Done()
	settings := telebot.Settings{
		Token:       os.Getenv("BOT_TOKEN"),
		Synchronous: true,
		Verbose:     true,
	}
	newTgBot, err := telebot.NewBot(settings)
	if err != nil {
		fmt.Println(err)
		panic("can't create bot")
	}
	tgBot = newTgBot
	tgBot.Handle(telebot.OnText, func(m *telebot.Message) {
		message := m.Text
		_, err := tgBot.Send(m.Sender, message)
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})
}

var dbClient *dynamodb.Client

func initDynamodb() {
	defer wg.Done()
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = os.Getenv("REGION")
		return nil
	})
	if err != nil {
		panic(err)
	}

	dbClient = dynamodb.NewFromConfig(cfg)
}

var adapter *httpadapter.HandlerAdapterV2

func initRouter() {
	defer wg.Done()
	r := httprouter.New()
	r.GET("/init-bot", initBot)
	r.POST("/bot", handleMessage)
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Not found", r.RequestURI)
		http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
	})
	adapter = httpadapter.NewV2(r)
}

var wg sync.WaitGroup

func init() {
	initSSM()

	wg = sync.WaitGroup{}
	wg.Add(3)

	go initTelebot()
	go initDynamodb()
	go initRouter()

	wg.Wait()
}

func initBot(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Set webhook
	wg.Add(2)

	go setWebhook()
	go setupDb()

	wg.Wait()
	returnOk(w)
}

func setWebhook() {
	defer wg.Done()
	err := tgBot.SetWebhook(&telebot.Webhook{
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: os.Getenv("domain") + "/" + os.Getenv("path_key") + "/bot",
		},
	})
	if err != nil {
		fmt.Println(err)
		panic("can't set webhook")
	}
}

type myRepo struct {
	PK            string `dynamodbav:"PK"`
	SK            string `dynamodbav:"SK"`
	GSI           string `dynamodbav:"GSI"`
	LSI           string `dynamodbav:"LSI"`
	Name          string `dynamodbav:"name"`
	Description   string `dynamodbav:"description"`
	AnyStingField string `dynamodbav:"anyStringField"`
	AnyIntField   int    `dynamodbav:"anyIntField"`
	AnyByteField  []byte `dynamodbav:"anyByteField"`
}

func listTables() ([]string, error) {
	var tableNames []string
	tables, err := dbClient.ListTables(
		context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Printf("Couldn't list tables. Here's why: %v\n", err)
	} else {
		tableNames = tables.TableNames
	}
	return tableNames, err
}

func createSubscribersTable() {
	defer wg.Done()
	_, err := dbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("name"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("url"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("name"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName:   aws.String(subscribersTable),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", subscribersTable, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(dbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(subscribersTable)}, 15*time.Second)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
	}
}

func createAdminsTable() {
	defer wg.Done()
	_, err := dbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("chatId"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("chatId"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName:   aws.String(adminsTable),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", adminsTable, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(dbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(adminsTable)}, 15*time.Second)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
	}
}

func contains(tables []string, table string) bool {
	for _, t := range tables {
		if t == table {
			return true
		}
	}
	return false
}

func setupDb() {
	defer wg.Done()
	tables, err := listTables()
	if err != nil {
		panic("can't list tables")
	}

	if !contains(tables, subscribersTable) {
		wg.Add(1)
		go createSubscribersTable()
	}

	if !contains(tables, adminsTable) {
		wg.Add(1)
		go createAdminsTable()
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Process update
	var u telebot.Update
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println(err)
		panic("can't unmarshal")
	}
	tgBot.ProcessUpdate(u)

	//returnOk(w)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(os.Getenv("domain") + "/" + os.Getenv("path_key") + "/bot"))
	if err != nil {
		fmt.Println(err)
		panic("can't write response")
	}
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
	handler := func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		// log path
		return adapter.ProxyWithContext(ctx, req)
		//// echo request
		//// get json from req
		//data, err := json.Marshal(req)
		//if err != nil {
		//	panic("can't marshal request")
		//}
		//
		//return events.APIGatewayV2HTTPResponse{
		//	Body:       string(data),
		//	StatusCode: 200,
		//}, nil

	}
	lambda.Start(handler)
}
