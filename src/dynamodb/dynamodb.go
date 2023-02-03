package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"os"
	"time"
	"website-checker-bot/src"
)

const (
	adminsTable      = "checker-admins"
	subscribersTable = "checker-subscribers"
)

var dbClient *dynamodb.Client

func initDynamodb() {
	defer main.Wg.Done()
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = os.Getenv("REGION")
		return nil
	})
	if err != nil {
		panic(err)
	}

	dbClient = dynamodb.NewFromConfig(cfg)
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
	defer main.Wg.Done()
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
	defer main.Wg.Done()
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
