package routes

import (
	"website-checker-bot/bot"
	"website-checker-bot/database"
	"website-checker-bot/ssm"
	"website-checker-bot/threadpool"
)

type Env struct {
	pool *threadpool.Pool
	bot  *bot.Bot
	db   *database.Db
}