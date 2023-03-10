package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (env *Env) InitBot(c *gin.Context) {
	// Set webhook
	url := os.Getenv("domain") + os.Getenv("path_key") + "/bot"
	env.pool.AddTask(func() {
		env.bot.SetWebhook(url)
	})
	env.pool.AddTask(env.db.Init)

	env.pool.Wait()

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
