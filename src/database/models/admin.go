package models

type Admin struct {
	ChatId  int64             `dynamodbav:"chat_id"`
	History map[string]string `dynamodbav:"history"`
}
