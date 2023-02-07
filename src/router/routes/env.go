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

func GetEnv() *Env {
	ssm.Init()

	pool := threadpool.New()
	tgBot := threadpool.MakeChan(bot.New)
	db := threadpool.MakeChan(func() *database.Db { return database.New(pool) })
	env := &Env{pool, <-tgBot, <-db}
	env.bot.Init(env)
	return env
}
