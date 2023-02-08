package admin

import (
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/bot/middlewares"

	"gopkg.in/telebot.v3"
)

func HandleGetAdmins(env *Env, c telebot.Context, args []string) error {
	if !CheckAdmin(env, c) {
		return nil
	}
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
