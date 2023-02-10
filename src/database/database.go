package database

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"os"
	"strconv"
	. "website-checker-bot/database/models"
	"website-checker-bot/threadpool"
)

const (
	AdminsTable   = "checker-admins"
	WebsitesTable = "checker-websites"
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

func (db *Db) Init() {}

func (db *Db) AddAdmin(chatId int64, username string) {
	data, err := attributevalue.MarshalMap(NewAdmin(chatId, username))
	if err != nil {
		log.Printf("Couldn't marshal admin. Here's why: %v\n", err)
	}
	_, err = db.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(AdminsTable),
		Item:      data,
	})
	if err != nil {
		log.Printf("Couldn't add admin. Here's why: %v\n", err)
	}
}

func (db *Db) GetAdmins() []Admin {
	var admins []Admin
	output, err := db.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(AdminsTable),
		AttributesToGet: []string{
			"chat_id",
			"username",
		},
	})
	if err != nil {
		log.Printf("Couldn't get admins. Here's why: %v\n", err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(output.Items, &admins)
		if err != nil {
			log.Printf("Couldn't unmarshal admins. Here's why: %v\n", err)
		}
		for i := range admins {
			if admins[i].History == nil {
				admins[i].History = make(map[string]string)
			}
		}
	}
	return admins
}

func (db *Db) GetAdmin(chatId int64) Admin {
	var admin Admin
	output, err := db.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(AdminsTable),
		Key: map[string]types.AttributeValue{
			"chat_id": &types.AttributeValueMemberN{Value: strconv.FormatInt(chatId, 10)},
		},
	})
	if err != nil {
		log.Printf("Couldn't get admin. Here's why: %v\n", err)
	} else {
		err = attributevalue.UnmarshalMap(output.Item, &admin)
		if err != nil {
			log.Printf("Couldn't unmarshal admin. Here's why: %v\n", err)
		}
		if admin.History == nil {
			admin.History = make(map[string]string)
		}
	}
	return admin
}

func (db *Db) UpdateAdmin(admin Admin) {
	data, err := attributevalue.MarshalMap(admin)
	if err != nil {
		log.Printf("Couldn't marshal admin. Here's why: %v\n", err)
	}
	_, err = db.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(AdminsTable),
		Item:      data,
	})
	if err != nil {
		log.Printf("Couldn't update admin. Here's why: %v\n", err)
	}
}

func (db *Db) DeleteAdmin(chatId int64) {
	_, err := db.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(AdminsTable),
		Key: map[string]types.AttributeValue{
			"chat_id": &types.AttributeValueMemberN{Value: strconv.FormatInt(chatId, 10)},
		},
	})
	if err != nil {
		log.Printf("Couldn't delete admin. Here's why: %v\n", err)
	}
}

func (db *Db) AddAdminMessage(chatId int64, part string, message string) {
	admin := db.GetAdmin(chatId)
	admin.History[part] = message
	db.UpdateAdmin(admin)
}

func (db *Db) AddWebsite(name string, url string, xpath string) {
	data, err := attributevalue.MarshalMap(NewWebsite(name, url, xpath))
	if err != nil {
		log.Printf("Couldn't marshal website. Here's why: %v\n", err)
	}
	_, err = db.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(WebsitesTable),
		Item:      data,
	})
	if err != nil {
		log.Printf("Couldn't add website. Here's why: %v\n", err)
	}
}

func (db *Db) CheckWebsite(name string) bool {
	output, err := db.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(WebsitesTable),
		Key: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		log.Printf("Couldn't check website. Here's why: %v\n", err)
		return false
	}
	return output.Item != nil
}

func (db *Db) CheckWebsiteUrl(url string) bool {
	output, err := db.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(WebsitesTable),
		Key: map[string]types.AttributeValue{
			"url": &types.AttributeValueMemberS{Value: url},
		},
	})
	if err != nil {
		log.Printf("Couldn't check website. Here's why: %v\n", err)
		return false
	}
	return output.Item != nil
}

func (db *Db) GetWebsites(withSubscribers bool) []Website {
	var websites []Website
	attr := []string{"name", "url", "xpath", "hash"}
	if withSubscribers {
		attr = append(attr, "subscribers")
	}
	output, err := db.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:       aws.String(WebsitesTable),
		AttributesToGet: attr,
	})

	if err != nil {
		log.Printf("Couldn't get websites. Here's why: %v\n", err)
		panic(err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(output.Items, &websites)
		if err != nil {
			log.Printf("Couldn't unmarshal websites. Here's why: %v\n", err)
			panic(err)
		}
		for i := range websites {
			if websites[i].Subscribers == nil {
				websites[i].Subscribers = make([]int64, 0)
			}
		}
	}
	return websites
}

func (db *Db) GetWebsite(name string) Website {
	var website Website
	output, err := db.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(WebsitesTable),
		Key: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		log.Printf("Couldn't get website. Here's why: %v\n", err)
	} else {
		err = attributevalue.UnmarshalMap(output.Item, &website)
		if err != nil {
			log.Printf("Couldn't unmarshal website. Here's why: %v\n", err)
		}
		if website.Subscribers == nil {
			website.Subscribers = make([]int64, 0)
		}
	}
	return website
}

func (db *Db) UpdateWebsite(website Website) {
	data, err := attributevalue.MarshalMap(website)
	if err != nil {
		log.Printf("Couldn't marshal website. Here's why: %v\n", err)
	}
	_, err = db.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(WebsitesTable),
		Item:      data,
	})
	if err != nil {
		log.Printf("Couldn't update website. Here's why: %v\n", err)
	}
}

func (db *Db) DeleteWebsite(name string) {
	_, err := db.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(WebsitesTable),
		Key: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		log.Printf("Couldn't delete website. Here's why: %v\n", err)
	}
}

func (db *Db) GetSubscriptions(id int64) []Website {
	websites := db.GetWebsites(true)
	var subscriptions []Website
	for _, website := range websites {
		for _, subscriber := range website.Subscribers {
			if subscriber == id {
				subscriptions = append(subscriptions, website)
			}
		}
	}
	return subscriptions
}

func (db *Db) AddSubscription(chatId int64, websiteName string) {
	website := db.GetWebsite(websiteName)
	website.Subscribers = append(website.Subscribers, chatId)
	db.UpdateWebsite(website)
}

func (db *Db) DeleteSubscription(chatId int64, websiteName string) {
	website := db.GetWebsite(websiteName)
	var subscribers []int64
	for _, subscriber := range website.Subscribers {
		if subscriber != chatId {
			subscribers = append(subscribers, subscriber)
		}
	}
	website.Subscribers = subscribers
	db.UpdateWebsite(website)
}

func (db *Db) CheckSubscription(chatId int64, websiteName string) bool {
	website := db.GetWebsite(websiteName)
	for _, subscriber := range website.Subscribers {
		if subscriber == chatId {
			return true
		}
	}
	return false
}
