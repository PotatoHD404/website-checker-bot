package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (env *Env) InitBot(c *gin.Context) {
	// Set webhook

	env.pool.AddTask(func() {
		env.bot.SetWebhook(os.Getenv("domain") + "/" + os.Getenv("path_key") + "/bot")
	})
	env.pool.AddTask(env.db.Init)

	env.pool.Wait()

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
