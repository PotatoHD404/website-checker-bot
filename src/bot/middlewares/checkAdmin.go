package middlewares

import (
	"fmt"
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/database/models"
	"website-checker-bot/utils"

	"gopkg.in/telebot.v3"
)

func CheckAdmin(env *Env, c telebot.Context) bool {
	if !utils.Contains(env.Db.GetAdmins(), NewAdmin(c.Sender().ID, c.Sender().Username)) {
		err := c.Reply("You are not admin")
		if err != nil {
			fmt.Println("Error checking website. Here is why: ", err)
			panic(err)
		}
		return false
	}
	return true
}
