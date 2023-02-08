package subscription

import (
	. "website-checker-bot/bot/commands/env"

	"gopkg.in/telebot.v3"
)

func HandleUnsubscribe(env *Env, c telebot.Context, args []string) error {
	// Process update
	//var u telebot.Update
	//
	//err := json.NewDecoder(c.Request.Body).Decode(&u)
	//if err != nil {
	//	fmt.Println(err)
	//	panic("can't unmarshal")
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "ok",
	//})
	return nil
}
