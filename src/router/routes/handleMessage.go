package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/tucnak/telebot.v2"
	"net/http"
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
