package admin

import (
	"strconv"
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/bot/middlewares"

	"gopkg.in/telebot.v3"
)

func HandleGetAdmins(env *Env, c telebot.Context, args []string) error {
	if !CheckAdmin(env, c) {
		return nil
	}
	admins := env.Db.GetAdmins()
	message := "Admins:\n"
	for _, admin := range admins {
		message += admin.Username + " - " + strconv.FormatInt(admin.ChatId, 10) + "\n"
	}
	err := c.Reply(message)
	if err != nil {
		return err
	}
	return nil
}
