package website

import (
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/bot/middlewares"

	"gopkg.in/telebot.v3"
)

func HandleDeleteWebsite(env *Env, c telebot.Context, args []string) error {
	if !CheckAdmin(env, c) {
		return nil
	}

	if len(args) != 1 {
		c.Reply("Usage: /delete_website <name>")
		return nil
	}

	websiteName := args[0]
	if !env.Db.CheckWebsite(websiteName) {
		err := c.Reply("Website with this name does not exist")
		if err != nil {
			return err
		}
		return nil
	}

	env.Db.DeleteWebsite(websiteName)

	err := c.Reply("Website deleted")
	if err != nil {
		return err
	}
	return nil
}
