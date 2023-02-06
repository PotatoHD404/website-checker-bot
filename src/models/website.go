package models

type Website struct {
	Name    string  `dynamodbav:"name"`
	Url     string  `dynamodbav:"url"`
	ChatIds []int64 `dynamodbav:"chat_ids"`
}
