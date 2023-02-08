package env

import (
	"website-checker-bot/database"
	"website-checker-bot/threadpool"
)

type Env struct {
	Pool *threadpool.Pool
	Db   *database.Db
}
