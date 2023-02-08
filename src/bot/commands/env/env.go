package env

import (
	"website-checker-bot/database"
	"website-checker-bot/threadpool"

	"gopkg.in/tucnak/telebot.v2"
)

type Env struct {
	Pool *threadpool.Pool
	Bot  *telebot.Bot
	Db   *database.Db
}
