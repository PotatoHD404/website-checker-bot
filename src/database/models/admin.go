package models

type Admin struct {
	ChatId   int64             `dynamodbav:"chat_id"`
	Username string            `dynamodbav:"username"`
	History  map[string]string `dynamodbav:"history"`
}

func NewAdmin(chatId int64, username string) Admin {
	return Admin{
		ChatId:   chatId,
		Username: username,
		History:  make(map[string]string),
	}
}
