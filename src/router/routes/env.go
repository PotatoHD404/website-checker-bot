package routes

import (
	"website-checker-bot/bot"
	b "website-checker-bot/bot/commands/env"
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
	tgBotCh := threadpool.MakeChan(bot.New)
	dbCh := threadpool.MakeChan(func() *database.Db { return database.New(pool) })
	tgBot, db := <-tgBotCh, <-dbCh
	env := &Env{pool, tgBot, db}
	env.bot.Init(&b.Env{Pool: pool, Db: db})

	return env
}
