package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/telebot.v3"
)

func (env *Env) HandleMessage(c *gin.Context) {
	// Process update
	var u telebot.Update

	err := json.NewDecoder(c.Request.Body).Decode(&u)
	if err != nil {
		fmt.Println(err)
		panic("can't unmarshal")
	}
	env.bot.ProcessUpdate(u)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
