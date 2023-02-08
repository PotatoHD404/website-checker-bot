package models

type Admin struct {
	ChatId  int64             `dynamodbav:"chat_id"`
	History map[string]string `dynamodbav:"history"`
}

func NewAdmin(chatId int64) Admin {
	return Admin{
		ChatId:  chatId,
		History: make(map[string]string),
	}
}
