package middlewares

import (
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/database/models"
	"website-checker-bot/utils"

	"gopkg.in/telebot.v3"
)

func CheckAdmin(env *Env, c telebot.Context) bool {
	if !utils.Contains(env.Db.GetAdmins(), NewAdmin(c.Sender().ID)) {
		c.Reply("You are not admin")
		return false
	}
	return true
}
