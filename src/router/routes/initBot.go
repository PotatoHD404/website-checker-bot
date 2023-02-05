package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/tucnak/telebot.v2"
	"net/http"
	"os"
)

func (env *Env) InitBot(c *gin.Context) {
	// Set webhook

	env.pool.AddTask(env.setWebhook)
	env.pool.AddTask(env.db.Init)

	env.pool.Wait()

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (env *Env) setWebhook() {
	err := env.bot.SetWebhook(&telebot.Webhook{
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: os.Getenv("domain") + "/" + os.Getenv("path_key") + "/bot",
		},
	})
	if err != nil {
		fmt.Println(err)
		panic("can't set webhook")
	}
}
