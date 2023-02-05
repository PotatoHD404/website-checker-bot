package database

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"os"
	"time"
	"website-checker-bot/threadpool"
	. "website-checker-bot/utils"
)

const (
	AdminsTable   = "checker-admins"
	WebsitesTable = "checker-subscribers"
)

type Db struct {
	pool   *threadpool.Pool
	client *dynamodb.Client
}

func New(pool *threadpool.Pool) *Db {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = os.Getenv("REGION")
		return nil
	})
	if err != nil {
		panic(err)
	}

	return &Db{pool, dynamodb.NewFromConfig(cfg)}
}

func (db *Db) Init() {
	tables, err := db.ListTables()
	if err != nil {
		panic("can't list tables")
	}

	if !Contains(tables, WebsitesTable) {
		db.pool.AddTask(db.CreateWebsitesTable)
	}

	if !Contains(tables, AdminsTable) {
		db.pool.AddTask(db.CreateAdminsTable)
	}
}

func (db *Db) ListTables() ([]string, error) {
	var tableNames []string
	tables, err := db.client.ListTables(
		context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Printf("Couldn't list tables. Here's why: %v\n", err)
	} else {
		tableNames = tables.TableNames
	}
	return tableNames, err
}

func (db *Db) CreateWebsitesTable() {
	_, err := db.client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("name"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("url"),
			AttributeType: types.ScalarAttributeTypeS,
		},
		//{
		//	AttributeName: aws.String("chatIds"),
		//	// array of strings
		//	AttributeType: types.ScalarAttributeTypeSS,
		//}
		},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("name"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName:   aws.String(WebsitesTable),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", WebsitesTable, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(db.client)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(WebsitesTable)}, 15*time.Second)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
	}
}

func (db *Db) CreateAdminsTable() {
	_, err := db.client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("chatId"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("chatId"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName:   aws.String(AdminsTable),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", AdminsTable, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(db.client)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(AdminsTable)}, 15*time.Second)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
	}
}
